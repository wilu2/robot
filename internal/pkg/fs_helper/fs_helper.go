package fshelper

import (
	"sort"
)

type TimeValue struct {
	Value     string
	Timestamp int64
}

func GetPeriodHeaderOrderIndex(headers []string) []string {

	sort.Slice(headers, func(j, k int) bool {
		return headers[j] < headers[k]
	})
	return headers
}
