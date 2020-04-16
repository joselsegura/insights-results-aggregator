#!/usr/bin/env python3

# Copyright © 2020 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Generate messages to be consumed by aggregator that are broken in some way.

This script read input message (that should be coorect) and generates bunch of
new messages. Each generated message is broken in some way so it is possible
to use such messages to test how broken messages are handled on aggregator
(ie. consumer) side.
"""

import json
import os
import sys
import copy
import itertools


def load_json(filename):
    """Load and decode JSON file."""
    with open(filename) as fin:
        return json.load(fin)


def generate_output(filename, payload):
    """Generate output JSON file with indentation."""
    with open(filename, 'w') as f:
        json.dump(payload, f, indent=4)
        print("Generated file {}".format(filename))


def filename_removed_items(removed_keys, selector=None):
    """Generate filename for JSON with items removed from original payload."""
    filename = "_".join(removed_keys)
    if selector is None:
        return "broken_without_attribute_" + filename
    else:
        return "broken_without_" + selector + "_attribute_" + filename


def remove_items_one_iter(original_payload, items_count, remove_flags,
                          selector=None):
    """One iteration of algorithm to remove items from original payload."""
    if selector is None:
        keys = list(original_payload.keys())
    else:
        keys = list(original_payload[selector].keys())

    # deep copy
    new_payload = copy.deepcopy(original_payload)
    removed_keys = []
    for i in range(items_count):
        remove_flag = remove_flags[i]
        if remove_flag:
            key = keys[i]
            if selector is None:
                del new_payload[key]
            else:
                del new_payload[selector][key]
            removed_keys.append(key)
    filename = filename_removed_items(removed_keys, selector)
    generate_output(filename, new_payload)


def remove_items(original_payload, selector=None):
    """Algorithm to remove items from original payload."""
    if selector is None:
        items_count = len(original_payload)
    else:
        items_count = len(original_payload[selector])

    # lexicographics ordering
    remove_flags_list = list(itertools.product([True, False],
                             repeat=items_count))
    # the last item contains (False, False, False...) and we are not interested
    # in removing ZERO items
    remove_flags_list = remove_flags_list[:-1]

    for remove_flags in remove_flags_list:
        remove_items_one_iter(original_payload, items_count, remove_flags,
                              selector)


def main(filename):
    """Entry point to this script."""
    original_payload = load_json(filename)
    remove_items(original_payload)
    remove_items(original_payload, "Report")


if len(sys.argv) < 2:
    print("Usage: python gen_broken_messages.py input_file.json")
    sys.exit(1)

main(sys.argv[1])
