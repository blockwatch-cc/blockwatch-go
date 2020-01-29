// Copyright (c) 2020 Blockwatch Data Inc.
// Author: alex@blockwatch.cc
//

package blockwatch

import (
	"fmt"
	"strings"
)

type OrderMode string

const (
	OrderInvalid OrderMode = ""
	OrderAsc     OrderMode = "asc"
	OrderDesc    OrderMode = "desc"
)

func ParseOrderModeIgnoreError(s string) OrderMode {
	m, _ := ParseOrderMode(s)
	return m
}

func ParseOrderMode(s string) (OrderMode, error) {
	switch strings.ToLower(s) {
	case "asc":
		return OrderAsc, nil
	case "desc", "":
		return OrderDesc, nil
	default:
		return OrderInvalid, fmt.Errorf("invalid order mode '%s'", s)
	}
}

func (m OrderMode) IsValid() bool {
	return m != OrderInvalid
}

func (m OrderMode) String() string {
	return string(m)
}

func (m OrderMode) MarshalText() ([]byte, error) {
	return []byte(m), nil
}

func (m *OrderMode) UnmarshalText(data []byte) error {
	mm, err := ParseOrderMode(string(data))
	if err != nil {
		return err
	}
	*m = mm
	return nil
}
