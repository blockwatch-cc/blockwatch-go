// Copyright (c) 2019 KIDTSUNAMI
// Author: alex@kidtsunami.com

package blockwatch

import (
	"fmt"
	"net/url"
	"strings"
)

type FilterMode int

const (
	FilterModeEqual FilterMode = iota
	FilterModeNotEqual
	FilterModeGt
	FilterModeGte
	FilterModeLt
	FilterModeLte
	FilterModeIn
	FilterModeNotIn
	FilterModeRange
	FilterModeRegexp
	FilterModeInvalid
)

type Filter struct {
	Field string
	Mode  FilterMode
	Value string
}

func NewFilter(field string, mode FilterMode, value string) *Filter {
	return &Filter{
		Field: field,
		Mode:  mode,
		Value: value,
	}
}

func (f Filter) String() string {
	return fmt.Sprintf("%s.%s=%s", f.Field, f.Mode, f.Value)
}

func (f Filter) AppendQuery(q url.Values) {
	q.Add(fmt.Sprintf("%s.%s", f.Field, f.Mode), f.Value)
}

func ParseFilterMode(s string) FilterMode {
	switch strings.ToLower(s) {
	case "", "eq":
		return FilterModeEqual
	case "ne":
		return FilterModeNotEqual
	case "gt":
		return FilterModeGt
	case "gte":
		return FilterModeGte
	case "lt":
		return FilterModeLt
	case "lte":
		return FilterModeLte
	case "in":
		return FilterModeIn
	case "nin":
		return FilterModeNotIn
	case "rg":
		return FilterModeRange
	case "re":
		return FilterModeRegexp
	default:
		return FilterModeInvalid
	}
}

func (m FilterMode) String() string {
	switch m {
	case FilterModeEqual:
		return "eq"
	case FilterModeNotEqual:
		return "ne"
	case FilterModeGt:
		return "gt"
	case FilterModeGte:
		return "gte"
	case FilterModeLt:
		return "lt"
	case FilterModeLte:
		return "lte"
	case FilterModeIn:
		return "in"
	case FilterModeNotIn:
		return "nin"
	case FilterModeRange:
		return "rg"
	case FilterModeRegexp:
		return "re"
	default:
		return "" // assuming equal
	}
}

// col_name.{ne|gt|gte|lt|lte|in|nin|re|rg}=value
func ParseFilter(key string, val string) (*Filter, error) {
	var fkey, mkey string
	if fields := strings.Split(key, "."); len(fields) == 2 {
		fkey, mkey = fields[0], fields[1]
	} else {
		fkey = fields[0]
		mkey = FilterModeEqual.String()
	}
	mode := ParseFilterMode(mkey)
	if mode == FilterModeInvalid {
		return nil, fmt.Errorf("invalid filter mode '%s'", mkey)
	}
	f := &Filter{
		Field: fkey,
		Mode:  mode,
		Value: val,
	}
	return f, nil
}
