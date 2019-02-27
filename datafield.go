// Copyright (c) 2019 KIDTSUNAMI
// Author: alex@kidtsunami.com

package blockwatch

import (
	"fmt"
	"strings"
)

type Datafield struct {
	Name string    `json:"name"`
	Code string    `json:"code"`
	Type FieldType `json:"type"`
}

type FieldType string

const (
	FieldTypeUndefined FieldType = ""
	FieldTypeString    FieldType = "string"
	FieldTypeBytes     FieldType = "bytes"
	FieldTypeDate      FieldType = "date"
	FieldTypeDatetime  FieldType = "datetime"
	FieldTypeBoolean   FieldType = "boolean"
	FieldTypeFloat64   FieldType = "float64"
	FieldTypeInt64     FieldType = "int64"
	FieldTypeUint64    FieldType = "uint64"
)

func ParseFieldType(s string) FieldType {
	switch strings.ToLower(s) {
	case "string":
		return FieldTypeString
	case "bytes":
		return FieldTypeBytes
	case "date":
		return FieldTypeDate
	case "datetime":
		return FieldTypeDatetime
	case "bool", "boolean":
		return FieldTypeBoolean
	case "integer", "int", "int64":
		return FieldTypeInt64
	case "unsigned", "uint", "uint64":
		return FieldTypeUint64
	case "float", "float64":
		return FieldTypeFloat64
	default:
		return FieldTypeUndefined
	}
}

func (f FieldType) String() string {
	switch f {
	case FieldTypeString:
		return "string"
	case FieldTypeBytes:
		return "bytes"
	case FieldTypeDate:
		return "date"
	case FieldTypeDatetime:
		return "datetime"
	case FieldTypeBoolean:
		return "boolean"
	case FieldTypeInt64:
		return "int64"
	case FieldTypeUint64:
		return "uint64"
	case FieldTypeFloat64:
		return "float64"
	default:
		return ""
	}
}

func (t FieldType) IsValid() bool {
	return t != FieldTypeUndefined
}

func (r FieldType) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

func (t *FieldType) UnmarshalText(data []byte) error {
	typ := ParseFieldType(string(data))
	if !typ.IsValid() {
		return fmt.Errorf("invalid datatfield ype %s", string(data))
	}
	*t = typ
	return nil
}
