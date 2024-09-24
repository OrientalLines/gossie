package utils

import (
	"fmt"
	"regexp"
	"strings"
)

type Tag struct {
	Name  string
	Value interface{}
}

var tagsRegex = regexp.MustCompile(`^gossie:(.+)$`)

func ParseTags(tags string) ([]Tag, error) {
	matches := tagsRegex.FindStringSubmatch(tags)
	if len(matches) != 2 {
		return nil, fmt.Errorf("invalid tags: %s", tags)
	}

	tagsSlice := strings.Split(matches[1], ",")
	result := make([]Tag, 0, len(tagsSlice))

	for _, tag := range tagsSlice {
		if i := strings.IndexByte(tag, '='); i >= 0 {
			result = append(result, Tag{
				Name:  tag[:i],
				Value: tag[i+1:],
			})
		} else {
			return nil, fmt.Errorf("invalid tag: %s", tag)
		}
	}

	return result, nil
}
