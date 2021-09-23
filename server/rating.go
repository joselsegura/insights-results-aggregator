// Copyright 2021 Red Hat, Inc
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

package server

import (
	"net/http"

	"github.com/RedHatInsights/insights-operator-utils/responses"
	"github.com/rs/zerolog/log"
)

func (server *HTTPServer) setRuleRating(writer http.ResponseWriter, request *http.Request) {
	rating, ok := readRuleRatingFromBody(writer, request)
	if !ok {
		// all errors handled inside
		return
	}

	userID, ok := readUserID(writer, request)
	if !ok {
		// everything is handled
		return
	}

	orgID, ok := readOrgID(writer, request)
	if !ok {
		// everything is handled
		return
	}

	ruleID, errorKey, err := getRuleAndErrorKeyFromRuleID(rating.Rule)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse rule identifier")
		handleServerError(writer, err)
		return
	}

	err = server.Storage.RateOnRule(userID, orgID, ruleID, errorKey, rating.Rating)
	if err != nil {
		log.Error().Err(err).Msg("Unable to store rating")
		handleServerError(writer, err)
		return
	}

	// If everythig goes fine, we should send the same ratings as response to the client
	err = responses.SendOK(writer, responses.BuildOkResponseWithData("ratings", rating))
	if err != nil {
		log.Error().Err(err).Msg("Errors sending response back to client")
		handleServerError(writer, err)
		return
	}
}
