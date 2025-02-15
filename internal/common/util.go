// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"errors"
	"strings"
)

// Contains checks if string contains any of the substrings.
func Contains(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// splitUserAndBotList splits the user and bot list from the request.
func splitUserAndBotList(ids []string) (botList []string, personList []string, err error) {
	for _, id := range ids {
		if strings.HasPrefix(id, "cli_") {
			botList = append(botList, id)
		} else if strings.HasPrefix(id, "ou_") {
			personList = append(personList, id)
		} else {
			return nil, nil, errors.New("invalid administrator ID")
		}
	}
	return botList, personList, nil
}
