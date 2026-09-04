package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Pointer Helper

func PtrString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func PtrInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

func PtrBool(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

func PtrUint64(p *uint64) uint64 {
	if p == nil {
		return 0
	}
	return *p
}

func PtrFloat64(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

func PtrStringSlice(s []*string) []string {
	if len(s) == 0 {
		return nil
	}
	out := make([]string, 0, len(s))
	for _, p := range s {
		if p == nil {
			continue
		}
		out = append(out, *p)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func ToUint64(p *string) uint64 {
	if p == nil {
		return 0
	}
	i, _ := strconv.ParseUint(*p, 10, 64)
	return i
}

// UnixSecondsToTime converts a Unix-seconds timestamp (*uint64) to *time.Time.
func UnixSecondsToTime(p *uint64) *time.Time {
	if p == nil || *p == 0 {
		return nil
	}
	t := time.Unix(int64(*p), 0).UTC()
	return &t
}

// UnixMillisToTime converts a Unix-milliseconds timestamp (*uint64) to *time.Time.
func UnixMillisToTime(p *uint64) *time.Time {
	if p == nil || *p == 0 {
		return nil
	}
	t := time.Unix(0, int64(*p)*int64(time.Millisecond)).UTC()
	return &t
}

// ParseTimestamp parses timestamp string into *time.Time.
func ParseTimestamp(value *string) *time.Time {
	if value == nil {
		return nil
	}

	raw := strings.TrimSpace(*value)
	if raw == "" {
		return nil
	}

	if isAllDigits(raw) {
		switch len(raw) {
		case 10: // Unix seconds
			if secs, err := strconv.ParseInt(raw, 10, 64); err == nil && secs > 0 {
				t := time.Unix(secs, 0).UTC()
				return &t
			}
		case 13: // Unix milliseconds
			if ms, err := strconv.ParseInt(raw, 10, 64); err == nil && ms > 0 {
				t := time.Unix(0, ms*int64(time.Millisecond)).UTC()
				return &t
			}
		}
	}

	for _, layout := range []string{
		time.DateTime,
		"2006-01-02 15:04:05.999999999",
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05+08:00",
		"2006-01-02 15:04:05-08:00",
		time.DateOnly,
	} {
		parsed, err := time.ParseInLocation(layout, raw, time.UTC)
		if err == nil {
			return &parsed
		}
	}

	return nil
}

// isAllDigits returns true if the string contains only ASCII digits.
func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// PageSizeFromLimit determines the page size based on the query's limit.
func PageSizeFromLimit(d *plugin.QueryData, defaultSize, minSize int64) int64 {
	if d.QueryContext == nil || d.QueryContext.Limit == nil {
		return defaultSize
	}
	limit := *d.QueryContext.Limit
	if limit >= defaultSize {
		return defaultSize
	}
	if limit < minSize {
		return minSize
	}
	return limit
}

// PageDone determines whether pagination is complete.
func PageDone[T int | int64 | uint64](offset, fetched int64, total T) bool {
	if fetched == 0 {
		return true
	}
	if total <= 0 {
		return false
	}
	return T(offset+fetched) >= total
}

// SingleIdFromQual extracts a single ID value from the query's equals-quals map.
func SingleIdFromQual(equalsQuals plugin.KeyColumnEqualsQualMap, qualName string) []*string {
	q := equalsQuals[qualName]
	if q == nil {
		return nil
	}
	v := q.GetStringValue()
	if v == "" {
		return nil
	}
	return []*string{&v}
}
