// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

// ConstructID constructs the ID for the resource.
// If the ID is empty, it will generate a random ID.
func ConstructID(resourceType TerraformType, resourceName TerraformName, id string) string {
	if id == "" {
		return string(resourceType) + "_" + string(resourceName) + "_" + time.Now().Format("20060102150405")
	}

	return string(resourceType) + "_" + string(resourceName) + "_" + id + "_" + time.Now().Format("20060102150405")
}

// StringValuesToStrings converts a list of basetypes.StringValue to a list of strings.
func StringValuesToStrings(values []basetypes.StringValue) []string {
	result := make([]string, 0, len(values))
	for _, v := range values {
		result = append(result, v.ValueString())
	}
	return result
}
