// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"errors"
	"strings"
)

// contains checks if string contains any of the substrings.
func contains(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// splitUserAndBotList splits the user and bot list from the request.
func splitUserAndBotList(request GroupChatAdministratorRequest) (botList []string, personList []string, err error) {
	for _, admin := range request.ManagerIDs {
		if strings.HasPrefix(admin, "cli_") {
			botList = append(botList, admin)
		} else if strings.HasPrefix(admin, "ou_") {
			personList = append(personList, admin)
		} else {
			return nil, nil, errors.New("invalid administrator ID")
		}
	}
	return botList, personList, nil
}
