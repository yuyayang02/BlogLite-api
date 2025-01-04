package adapter

import "strings"

func parseTagsFromTagParams(tags *string) []string {
	if tags != nil && *tags != "" {
		return strings.Split(*tags, ",")
	}
	return nil
}
