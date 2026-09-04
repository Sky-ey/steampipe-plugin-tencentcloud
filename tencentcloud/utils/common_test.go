package utils

import (
	"reflect"
	"testing"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func ptr[T any](v T) *T { return &v }

func TestPtrString(t *testing.T) {
	if got := PtrString(nil); got != "" {
		t.Fatalf("expected empty string for nil, got %q", got)
	}
	if got := PtrString(ptr("hello")); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestPtrInt64(t *testing.T) {
	if got := PtrInt64(nil); got != 0 {
		t.Fatalf("expected 0 for nil, got %d", got)
	}
	if got := PtrInt64(ptr(int64(42))); got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
	if got := PtrInt64(ptr(int64(-7))); got != -7 {
		t.Fatalf("expected -7, got %d", got)
	}
}

func TestPtrBool(t *testing.T) {
	if got := PtrBool(nil); got != false {
		t.Fatalf("expected false for nil, got %v", got)
	}
	if got := PtrBool(ptr(true)); got != true {
		t.Fatalf("expected true, got %v", got)
	}
	if got := PtrBool(ptr(false)); got != false {
		t.Fatalf("expected false, got %v", got)
	}
}

func TestPtrUint64(t *testing.T) {
	if got := PtrUint64(nil); got != 0 {
		t.Fatalf("expected 0 for nil, got %d", got)
	}
	if got := PtrUint64(ptr(uint64(100))); got != 100 {
		t.Fatalf("expected 100, got %d", got)
	}
	// value larger than int32 range to confirm it is not truncated
	if got := PtrUint64(ptr(uint64(1) << 40)); got != uint64(1)<<40 {
		t.Fatalf("expected %d, got %d", uint64(1)<<40, got)
	}
}

func TestPtrFloat64(t *testing.T) {
	if got := PtrFloat64(nil); got != 0 {
		t.Fatalf("expected 0 for nil, got %v", got)
	}
	if got := PtrFloat64(ptr(3.14)); got != 3.14 {
		t.Fatalf("expected 3.14, got %v", got)
	}
}

func TestPtrStringSlice(t *testing.T) {
	tests := []struct {
		name  string
		input []*string
		want  []string
	}{
		{
			name:  "nil input",
			input: nil,
			want:  nil,
		},
		{
			name:  "empty input",
			input: []*string{},
			want:  nil,
		},
		{
			name:  "all nil elements",
			input: []*string{nil, nil},
			want:  nil,
		},
		{
			name:  "skips nil elements and keeps order",
			input: []*string{ptr("a"), nil, ptr("b")},
			want:  []string{"a", "b"},
		},
		{
			name:  "all valid elements",
			input: []*string{ptr("x"), ptr("y"), ptr("z")},
			want:  []string{"x", "y", "z"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PtrStringSlice(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestParseTimestamp(t *testing.T) {
	timePtr := func(year int, month time.Month, day, hour, minute, second, nanosecond int) *time.Time {
		value := time.Date(year, month, day, hour, minute, second, nanosecond, time.UTC)
		return &value
	}

	tests := []struct {
		name  string
		value *string
		want  *time.Time
	}{
		{name: "nil", value: nil, want: nil},
		{name: "blank", value: ptr("   "), want: nil},
		{name: "invalid", value: ptr("not-a-time"), want: nil},
		{name: "rfc3339", value: ptr("2026-07-29T12:34:56Z"), want: timePtr(2026, 7, 29, 12, 34, 56, 0)},
		{name: "space separated", value: ptr("2026-07-29 12:34:56"), want: timePtr(2026, 7, 29, 12, 34, 56, 0)},
		{name: "fractional seconds", value: ptr("2026-07-29 12:34:56.123456"), want: timePtr(2026, 7, 29, 12, 34, 56, 123456000)},
		{name: "date", value: ptr("2026-07-29"), want: timePtr(2026, 7, 29, 0, 0, 0, 0)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ParseTimestamp(test.value)
			if test.want == nil {
				if got != nil {
					t.Fatalf("ParseTimestamp() = %v, want nil", got)
				}
				return
			}
			if got == nil || !got.Equal(*test.want) {
				t.Fatalf("ParseTimestamp() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestPageDone(t *testing.T) {
	tests := []struct {
		name            string
		offset, fetched int64
		total           int64
		want            bool
	}{
		{name: "empty page", fetched: 0, total: 100, want: true},
		{name: "last page reached", offset: 90, fetched: 10, total: 100, want: true},
		{name: "more pages", offset: 50, fetched: 50, total: 100, want: true},
		{name: "not done", offset: 30, fetched: 50, total: 100, want: false},
		{name: "total missing continues", offset: 0, fetched: 100, total: 0, want: false},
		{name: "negative total continues", offset: 0, fetched: 1, total: -1, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := PageDone(test.offset, test.fetched, test.total); got != test.want {
				t.Fatalf("PageDone(%d, %d, %d) = %v, want %v", test.offset, test.fetched, test.total, got, test.want)
			}
		})
	}
}

func TestPageSizeFromLimit(t *testing.T) {
	newQueryData := func(limit *int64) *plugin.QueryData {
		d := &plugin.QueryData{}
		if limit != nil {
			d.QueryContext = &plugin.QueryContext{Limit: limit}
		}
		return d
	}

	tests := []struct {
		name        string
		limit       *int64
		defaultSize int64
		minSize     int64
		want        int64
	}{
		{name: "no limit uses default", limit: nil, defaultSize: 100, minSize: 1, want: 100},
		{name: "nil query context uses default", limit: nil, defaultSize: 100, minSize: 1, want: 100},
		{name: "limit >= default uses default", limit: ptr(int64(500)), defaultSize: 100, minSize: 1, want: 100},
		{name: "limit equal default", limit: ptr(int64(100)), defaultSize: 100, minSize: 1, want: 100},
		{name: "limit below default converges", limit: ptr(int64(5)), defaultSize: 100, minSize: 1, want: 5},
		{name: "limit below minSize clamps", limit: ptr(int64(0)), defaultSize: 100, minSize: 1, want: 1},
		{name: "minSize larger than limit", limit: ptr(int64(3)), defaultSize: 100, minSize: 5, want: 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := newQueryData(test.limit)
			if got := PageSizeFromLimit(d, test.defaultSize, test.minSize); got != test.want {
				t.Fatalf("PageSizeFromLimit() = %d, want %d", got, test.want)
			}
		})
	}
}
