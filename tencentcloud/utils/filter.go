package utils

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// FilterMapping maps a steampipe qual name to a cloud API filter name.
type FilterMapping struct {
	QualName   string
	FilterName string
}

// BuildFilters builds a slice of SDK Filter objects from the query's equals-quals map according to the provided mappings.
func BuildFilters[T any](
	equalsQuals plugin.KeyColumnEqualsQualMap,
	newFilter func(name, value string) T,
	mappings ...FilterMapping,
) []T {
	filters := make([]T, 0, len(mappings))
	for _, m := range mappings {
		qual := equalsQuals[m.QualName]
		if qual == nil {
			continue
		}
		value := qual.GetStringValue()
		if value == "" {
			continue
		}
		filters = append(filters, newFilter(m.FilterName, value))
	}
	return filters
}
