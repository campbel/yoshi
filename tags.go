package yoshi

import (
	"reflect"
	"strings"
)

const (
	tagFlag = "yoshi"
)

type Tag struct {
	Flags       []string
	Description string
	Default     string
}

func getTags(field reflect.StructField) Tag {
	return parseTag(field.Tag.Get(tagFlag))
}

func parseTag(tag string) Tag {
	if tag == "" {
		return Tag{}
	}
	parts := strings.Split(tag, ";")
	if len(parts) == 0 {
		return Tag{}
	}
	if len(parts) == 1 {
		return Tag{strings.Split(parts[0], ","), "", ""}
	}
	if len(parts) == 2 {
		return Tag{strings.Split(parts[0], ","), parts[1], ""}
	}
	return Tag{strings.Split(parts[0], ","), parts[1], parts[2]}
}
