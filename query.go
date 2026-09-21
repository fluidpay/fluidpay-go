package fluidpay

import "time"

// Search operators accepted by the gateway.
const (
	OpEqual       = "="
	OpNotEqual    = "!="
	OpLessThan    = "<"
	OpGreaterThan = ">"
)

// StringQuery matches a string field. Supported operators are "=" and "!=".
type StringQuery struct {
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// IntQuery matches an integer field. Supported operators are "=", "!=", "<"
// and ">".
type IntQuery struct {
	Operator string `json:"operator"`
	Value    int    `json:"value"`
}

// DateRangeQuery matches a timestamp between StartDate and EndDate. Dates are
// UTC in the form "2006-01-02T15:04:05Z".
type DateRangeQuery struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// Equal builds a StringQuery matching value exactly.
func Equal(value string) *StringQuery { return &StringQuery{Operator: OpEqual, Value: value} }

// NotEqual builds a StringQuery excluding value.
func NotEqual(value string) *StringQuery { return &StringQuery{Operator: OpNotEqual, Value: value} }

// IntEqual builds an IntQuery matching value exactly.
func IntEqual(value int) *IntQuery { return &IntQuery{Operator: OpEqual, Value: value} }

// IntNotEqual builds an IntQuery excluding value.
func IntNotEqual(value int) *IntQuery { return &IntQuery{Operator: OpNotEqual, Value: value} }

// LessThan builds an IntQuery matching values below value.
func LessThan(value int) *IntQuery { return &IntQuery{Operator: OpLessThan, Value: value} }

// GreaterThan builds an IntQuery matching values above value.
func GreaterThan(value int) *IntQuery { return &IntQuery{Operator: OpGreaterThan, Value: value} }

// DateRange builds a DateRangeQuery covering [start, end]. Times are
// converted to UTC and formatted the way the gateway expects.
func DateRange(start, end time.Time) *DateRangeQuery {
	return &DateRangeQuery{
		StartDate: start.UTC().Format(time.RFC3339),
		EndDate:   end.UTC().Format(time.RFC3339),
	}
}

// Day builds a DateRangeQuery covering a single calendar day in UTC.
func Day(day time.Time) *DateRangeQuery {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	return DateRange(start, start.Add(24*time.Hour-time.Second))
}
