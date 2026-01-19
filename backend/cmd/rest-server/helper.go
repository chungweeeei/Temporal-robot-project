package main

import (
	"strconv"
	"strings"

	"go.temporal.io/sdk/client"
)

type ReadableCalendarSpec struct {
	Second  string `json:"second,omitempty"`
	Minute  string `json:"minute,omitempty"`
	Hour    string `json:"hour,omitempty"`
	Month   string `json:"month,omitempty"`
	Year    string `json:"year,omitempty"`
	Comment string `json:"comment,omitempty"`
}

func rangeToString(ranges []client.ScheduleRange) string {

	if len(ranges) == 0 {
		return "*"
	}

	var parts []string
	for _, r := range ranges {
		s := strconv.Itoa(r.Start)
		if r.End > r.Start {
			s += "-" + strconv.Itoa(r.End)
		}
		if r.Step > 1 {
			s += "/" + strconv.Itoa(r.Step)
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, ",")
}
