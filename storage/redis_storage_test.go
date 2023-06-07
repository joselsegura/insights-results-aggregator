// Copyright 2023 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage_test

import (
	"fmt"
	"time"

	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/insights-operator-utils/redis"

	"github.com/RedHatInsights/insights-results-aggregator-data/testdata"
	"github.com/RedHatInsights/insights-results-aggregator/storage"
)

// default Redis configuration
var configuration = storage.Configuration{
	RedisConfiguration: storage.RedisConfiguration{
		RedisEndpoint:       "localhost:12345",
		RedisDatabase:       0,
		RedisTimeoutSeconds: 1,
		RedisPassword:       "",
	},
}

// getMockRedis is used to get a mocked Redis client to expect and
// respond to queries
func getMockRedis(t *testing.T) (
	mockClient storage.RedisStorage, mockServer redismock.ClientMock,
) {
	client, mockServer := redismock.NewClientMock()
	mockClient = storage.RedisStorage{
		Client: redis.Client{Connection: client},
	}
	err := mockClient.Init()
	if err != nil {
		t.Fatal(err)
	}
	return
}

// assertRedisExpectationsMet helper function used to ensure mock expectations were met
func assertRedisExpectationsMet(t *testing.T, mock redismock.ClientMock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

// TestNewRedisClient checks if it is possible to construct Redis client
func TestNewRedisClient(t *testing.T) {
	// try to instantiate Redis storage
	client, err := storage.NewRedisStorage(configuration)

	// check results
	assert.NotNil(t, client)
	assert.NoError(t, err)
}

// TestNewDummyRedisClient checks if it is possible to construct Redis
// client structure useful for testing
func TestNewDummyRedisClient(t *testing.T) {
	// configuration where Redis endpoint is set to empty string
	configuration := storage.Configuration{
		RedisConfiguration: storage.RedisConfiguration{
			RedisEndpoint:       "",
			RedisDatabase:       0,
			RedisTimeoutSeconds: 1,
			RedisPassword:       "",
		},
	}
	// try to instantiate Redis storage
	client, err := storage.NewRedisStorage(configuration)

	// check results
	assert.NotNil(t, client)
	assert.NoError(t, err)
}

// TestNewRedisClientDBIndexOutOfRange checks if Redis client
// constructor checks for incorrect database index
func TestNewRedisClientDBIndexOutOfRange(t *testing.T) {
	configuration1 := storage.Configuration{
		RedisConfiguration: storage.RedisConfiguration{
			RedisEndpoint:       "localhost:12345",
			RedisDatabase:       -1,
			RedisTimeoutSeconds: 1,
			RedisPassword:       "",
		},
	}
	// try to instantiate Redis storage
	client, err := storage.NewRedisStorage(configuration1)

	// check results
	assert.Nil(t, client)
	assert.Error(t, err)

	configuration2 := storage.Configuration{
		RedisConfiguration: storage.RedisConfiguration{
			RedisEndpoint:       "localhost:12345",
			RedisDatabase:       16,
			RedisTimeoutSeconds: 1,
			RedisPassword:       "",
		},
	}
	// try to instantiate Redis storage
	client, err = storage.NewRedisStorage(configuration2)

	// check results
	assert.Nil(t, client)
	assert.Error(t, err)
}

// TestWrittenReportsMetric checks the method WriteReportForCluster
func TestRediWriteReportForCluster(t *testing.T) {
	client, server := getMockRedis(t)
	expectedKey := fmt.Sprintf("organization:%d:cluster:%s:request:%s",
		int(testdata.OrgID),
		string(testdata.ClusterName),
		string(testdata.RequestID1))

	// it is expected that key will be set with given expiration period
	server.ExpectSet(expectedKey, "", client.Expiration).SetVal("OK")

	err := client.WriteReportForCluster(
		testdata.OrgID, testdata.ClusterName,
		testdata.Report3Rules, testdata.Report3RulesParsed,
		testdata.LastCheckedAt, testdata.LastCheckedAt, time.Now(),
		testdata.KafkaOffset, testdata.RequestID1)

	assert.NoError(t, err)
	assertRedisExpectationsMet(t, server)
}
