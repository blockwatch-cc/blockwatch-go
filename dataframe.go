// Copyright (c) 2019 KIDTSUNAMI
// Author: alex@kidtsunami.com

package blockwatch

import (
	"bytes"
	"encoding"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Dataframe struct {
	Columns []Datafield       `json:"columns"`
	Data    []json.RawMessage `json:"data"`

	// internal fields used for data access
	tinfo  *typeInfo
	colmap map[string]int
}

type Row struct {
	data *Dataframe
	n    int
}

func (r Row) Decode(val interface{}) error {
	return r.data.DecodeAt(r.n, val)
}

func (r Row) Dataframe() *Dataframe {
	return r.data
}

func (r Row) Column(name string) (int, interface{}, error) {
	col := r.data.columnIndex(name)
	if col < 0 {
		return -1, nil, fmt.Errorf("blockwatch: missing column '%s'", name)
	}
	v, err := r.data.FieldAt(col, r.n)
	return col, v, err
}

func (t *Dataframe) ResetType() {
	t.tinfo = nil
	t.colmap = nil
}

func (t *Dataframe) DecodeAt(row int, val interface{}) error {
	if t.tinfo == nil {
		var err error
		t.tinfo, err = getTypeInfo(val)
		if err != nil {
			return err
		}
		t.colmap = make(map[string]int, len(t.Columns))
		for i, v := range t.Columns {
			t.colmap[v.Code] = i
		}
	}
	return t.decodeAt(row, val, t.tinfo)
}

func (t *Dataframe) FieldAt(col, row int) (interface{}, error) {
	if len(t.Columns) < col {
		return nil, fmt.Errorf("blockwatch: invalid data column %d > len %d", col, len(t.Columns))
	}
	if len(t.Data) < row {
		return nil, fmt.Errorf("blockwatch: invalid data row %d > len %d", row, len(t.Data))
	}
	name := t.Columns[col].Code
	switch typ := t.Columns[col].Type; typ {
	case FieldTypeString:
		return t.decodeStringAt(col, row, name)
	case FieldTypeBytes:
		return t.decodeBytesAt(col, row, name)
	case FieldTypeDate, FieldTypeDatetime:
		v, err := t.decodeTimeAt(col, row, name)
		return v, err
	case FieldTypeBoolean:
		return t.decodeBoolAt(col, row, name)
	case FieldTypeFloat64:
		return t.decodeFloat64At(col, row, name)
	case FieldTypeInt64:
		return t.decodeInt64At(col, row, name)
	case FieldTypeUint64:
		return t.decodeUint64At(col, row, name)
	default:
		return nil, fmt.Errorf("blockwatch: no method for decoding column '%s' type %s",
			name, typ)
	}
}

func (t *Dataframe) ForEach(fn func(r Row) error) error {
	for i, l := 0, len(t.Data); i < l; i++ {
		if err := fn(Row{data: t, n: i}); err != nil {
			return err
		}
	}
	return nil
}

func (t *Dataframe) Column(name string) (int, interface{}, error) {
	i := t.columnIndex(name)
	if i < 0 {
		return -1, nil, fmt.Errorf("blockwatch: missing column '%s'", name)
	}
	switch t.Columns[i].Type {
	case FieldTypeString:
		v, err := t.decodeStringColumn(i, name)
		return i, v, err
	case FieldTypeBytes:
		v, err := t.decodeBytesColumn(i, name)
		return i, v, err
	case FieldTypeDate, FieldTypeDatetime:
		v, err := t.decodeTimeColumn(i, name)
		return i, v, err
	case FieldTypeBoolean:
		v, err := t.decodeBoolColumn(i, name)
		return i, v, err
	case FieldTypeFloat64:
		v, err := t.decodeFloat64Column(i, name)
		return i, v, err
	case FieldTypeInt64:
		v, err := t.decodeInt64Column(i, name)
		return i, v, err
	case FieldTypeUint64:
		v, err := t.decodeUint64Column(i, name)
		return i, v, err
	default:
		return i, nil, fmt.Errorf("blockwatch: no method for decoding column '%s' type %s",
			name, t.Columns[i].Type)
	}
}

func (t *Dataframe) columnIndex(name string) int {
	if i, ok := t.colmap[name]; ok {
		return i
	}
	return -1
}

func (t *Dataframe) decodeAt(pos int, val interface{}, tinfo *typeInfo) error {
	if len(t.Data) <= pos {
		return fmt.Errorf("blockwatch: invalid table row %d > len %d", pos, len(t.Data))
	}
	v := derefValue(reflect.ValueOf(val))
	if !v.IsValid() {
		return fmt.Errorf("blockwatch: invalid value of type %T", v)
	}
	// decode from json.RawMessage
	buf := t.Data[pos]

	// strip JSON array delimiter and split JSON columns
	cols := bytes.Split(buf[1:len(buf)-1], []byte(","))

	// decode all struct fields
	for _, finfo := range tinfo.fields {
		// determine column related to struct field based on name, skip missing struct fields
		col := t.columnIndex(finfo.name)
		if col < 0 {
			continue
		}
		// resolve field, fail on error
		dst := finfo.value(v)
		if !dst.IsValid() {
			return fmt.Errorf("blockwatch: invalid struct field value for field %s/%s[%s]", finfo.name, dst.Type().String(), dst.Kind().String())
		}
		dst0 := dst
		// deref pointers and allocate if nil
		if dst.Kind() == reflect.Ptr {
			if dst.IsNil() && dst.CanSet() {
				dst.Set(reflect.New(dst.Type().Elem()))
			}
			dst = dst.Elem()
		}

		// unless dst is time.Time we call binary and text unamrshalers for custom types
		if dst.Type().String() != "time.Time" {
			// try binary unmarshalers first
			if dst.CanInterface() && dst.Type().Implements(binaryUnmarshalerType) {
				if err := dst.Interface().(encoding.BinaryUnmarshaler).UnmarshalBinary(cols[col]); err != nil {
					return makeFieldError(finfo.name, dst, err)
				}
				continue
			}

			if dst.CanAddr() {
				pv := dst.Addr()
				if pv.CanInterface() && pv.Type().Implements(binaryUnmarshalerType) {
					if err := pv.Interface().(encoding.BinaryUnmarshaler).UnmarshalBinary(cols[col]); err != nil {
						return makeFieldError(finfo.name, dst, err)
					}
					continue
				}
			}

			// try text unmarshalers next
			if dst.CanInterface() && dst.Type().Implements(textUnmarshalerType) {
				if err := dst.Interface().(encoding.TextUnmarshaler).UnmarshalText(cols[col]); err != nil {
					return makeFieldError(finfo.name, dst, err)
				}
				continue
			}

			if dst.CanAddr() {
				pv := dst.Addr()
				if pv.CanInterface() && pv.Type().Implements(textUnmarshalerType) {
					if err := pv.Interface().(encoding.TextUnmarshaler).UnmarshalText(cols[col]); err != nil {
						return makeFieldError(finfo.name, dst, err)
					}
					continue
				}
			}
		}
		// unmarshal simple values
		colstr := string(cols[col])
		switch dst.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.ParseInt(colstr, 10, dst.Type().Bits())
			if err != nil {
				return makeFieldError(finfo.name, dst, err)
			}
			dst.SetInt(i)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			i, err := strconv.ParseUint(colstr, 10, dst.Type().Bits())
			if err != nil {
				return makeFieldError(finfo.name, dst, err)
			}
			dst.SetUint(i)
		case reflect.Float32, reflect.Float64:
			i, err := strconv.ParseFloat(colstr, dst.Type().Bits())
			if err != nil {
				return makeFieldError(finfo.name, dst, err)
			}
			dst.SetFloat(i)
		case reflect.Bool:
			i, err := strconv.ParseBool(strings.TrimSpace(colstr))
			if err != nil {
				return makeFieldError(finfo.name, dst, err)
			}
			dst.SetBool(i)
		case reflect.String:
			i, err := strconv.Unquote(strings.TrimSpace(colstr))
			if err != nil {
				return makeFieldError(finfo.name, dst, err)
			}
			dst.SetString(i)
		case reflect.Slice:
			// make sure it's a byte slice
			if dst.Type().Elem().Kind() == reflect.Uint8 {
				str, err := strconv.Unquote(strings.TrimSpace(colstr))
				if err != nil {
					return makeFieldError(finfo.name, dst, err)
				}
				buf, err := hex.DecodeString(str)
				if err != nil {
					return makeFieldError(finfo.name, dst, err)
				}
				dst.SetBytes(buf)
			} else {
				return fmt.Errorf("blockwatch: unsupported embedded slice type %s", dst.Type().Elem().Kind().String())
			}
		case reflect.Struct:
			// special time.Time decoding (Blockwatch JSON contains UNIX millisec)
			if dst.Type().String() == "time.Time" {
				i, err := strconv.ParseInt(colstr, 10, 64)
				if err != nil {
					return err
				}
				tv := reflect.ValueOf(time.Unix(0, i*1000000).UTC())
				dst.Set(tv)
			} else {
				return fmt.Errorf("blockwatch: unsupported embedded struct type %s", dst.Type().String())
			}
		default:
			return fmt.Errorf("blockwatch: no method for unmarshaling type %s (%s)", dst0.Type().String(), dst0.Kind().String())
		}
	}
	return nil
}

func (t *Dataframe) decodeInt64Column(col int, name string) ([]int64, error) {
	vec := make([]int64, len(t.Data))
	var err error
	for i, _ := range t.Data {
		vec[i], err = t.decodeInt64At(col, i, name)
		if err != nil {
			return nil, err
		}
	}
	return vec, nil
}

func (t *Dataframe) decodeUint64Column(col int, name string) ([]uint64, error) {
	vec := make([]uint64, len(t.Data))
	var err error
	for i, _ := range t.Data {
		vec[i], err = t.decodeUint64At(col, i, name)
		if err != nil {
			return nil, err
		}
	}
	return vec, nil
}

func (t *Dataframe) decodeFloat64Column(col int, name string) ([]float64, error) {
	vec := make([]float64, len(t.Data))
	var err error
	for i, _ := range t.Data {
		vec[i], err = t.decodeFloat64At(col, i, name)
		if err != nil {
			return nil, err
		}
	}
	return vec, nil
}

func (t *Dataframe) decodeStringColumn(col int, name string) ([]string, error) {
	vec := make([]string, len(t.Data))
	var err error
	for i, _ := range t.Data {
		vec[i], err = t.decodeStringAt(col, i, name)
		if err != nil {
			return nil, err
		}
	}
	return vec, nil
}

func (t *Dataframe) decodeBytesColumn(col int, name string) ([][]byte, error) {
	vec := make([][]byte, len(t.Data))
	var err error
	for i, _ := range t.Data {
		vec[i], err = t.decodeBytesAt(col, i, name)
		if err != nil {
			return nil, err
		}
	}
	return vec, nil
}

func (t *Dataframe) decodeBoolColumn(col int, name string) ([]bool, error) {
	vec := make([]bool, len(t.Data))
	var err error
	for i, _ := range t.Data {
		vec[i], err = t.decodeBoolAt(col, i, name)
		if err != nil {
			return nil, err
		}
	}
	return vec, nil
}

func (t *Dataframe) decodeTimeColumn(col int, name string) ([]time.Time, error) {
	vec := make([]time.Time, len(t.Data))
	var err error
	for i, _ := range t.Data {
		vec[i], err = t.decodeTimeAt(col, i, name)
		if err != nil {
			return nil, err
		}
	}
	return vec, nil
}

func (t *Dataframe) decodeInt64At(col, row int, name string) (int64, error) {
	// strip JSON array delimiters
	v := t.Data[row]
	v = v[1 : len(v)-1]

	// find the n-th column separated by comma
	start, end := indexByteColumnN(v, ',', col)
	if start < 0 {
		return 0, makeColumnMissingError(name, col, row)
	}
	val, err := strconv.ParseInt(string(v[start:end]), 10, 64)
	if err != nil {
		return 0, makeColumnError(name, col, row, err)
	}
	return val, nil
}

func (t *Dataframe) decodeUint64At(col, row int, name string) (uint64, error) {
	// strip JSON array delimiters
	v := t.Data[row]
	v = v[1 : len(v)-1]

	// find the n-th column separated by comma
	start, end := indexByteColumnN(v, ',', col)
	if start < 0 {
		return 0, makeColumnMissingError(name, col, row)
	}
	val, err := strconv.ParseUint(string(v[start:end]), 10, 64)
	if err != nil {
		return 0, makeColumnError(name, col, row, err)
	}
	return val, nil
}

func (t *Dataframe) decodeFloat64At(col, row int, name string) (float64, error) {
	// strip JSON array delimiters
	v := t.Data[row]
	v = v[1 : len(v)-1]

	// find the n-th column separated by comma
	start, end := indexByteColumnN(v, ',', col)
	if start < 0 {
		return 0, makeColumnMissingError(name, col, row)
	}
	val, err := strconv.ParseFloat(string(v[start:end]), 64)
	if err != nil {
		return 0, makeColumnError(name, col, row, err)
	}
	return val, nil
}

func (t *Dataframe) decodeStringAt(col, row int, name string) (string, error) {
	// strip JSON array delimiters
	v := t.Data[row]
	v = v[1 : len(v)-1]

	// find the n-th column separated by comma
	start, end := indexByteColumnN(v, ',', col)
	if start < 0 {
		return "", makeColumnMissingError(name, col, row)
	}
	val, err := strconv.Unquote(string(bytes.TrimSpace(v[start:end])))
	if err != nil {
		return "", makeColumnError(name, col, row, err)
	}
	return val, nil
}

func (t *Dataframe) decodeBytesAt(col, row int, name string) ([]byte, error) {
	// strip JSON array delimiters
	v := t.Data[row]
	v = v[1 : len(v)-1]

	// find the n-th column separated by comma
	start, end := indexByteColumnN(v, ',', col)
	if start < 0 {
		return nil, makeColumnMissingError(name, col, row)
	}
	val, err := strconv.Unquote(string(bytes.TrimSpace(v[start:end])))
	if err != nil {
		return nil, makeColumnError(name, col, row, err)
	}
	buf, err := hex.DecodeString(val)
	if err != nil {
		return nil, makeColumnError(name, col, row, err)
	}
	return buf, nil
}

func (t *Dataframe) decodeBoolAt(col, row int, name string) (bool, error) {
	// strip JSON array delimiters
	v := t.Data[row]
	v = v[1 : len(v)-1]

	// find the n-th column separated by comma
	start, end := indexByteColumnN(v, ',', col)
	if start < 0 {
		return false, makeColumnMissingError(name, col, row)
	}
	val, err := strconv.ParseBool(string(v[start:end]))
	if err != nil {
		return false, makeColumnError(name, col, row, err)
	}
	return val, nil
}

func (t *Dataframe) decodeTimeAt(col, row int, name string) (time.Time, error) {
	// strip JSON array delimiters
	v := t.Data[row]
	v = v[1 : len(v)-1]

	// find the n-th column separated by comma
	start, end := indexByteColumnN(v, ',', col)
	if start < 0 {
		return time.Time{}, fmt.Errorf("blockwatch: missing column %s (%d:%d)", name, col, row)
	}
	val, err := strconv.ParseInt(string(v[start:end]), 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("blockwatch: cannot decode column %s [%d:%d]: %v",
			name, col, row, err)
	}
	return time.Unix(0, val*1000000).UTC(), nil
}

func makeFieldError(name string, val reflect.Value, err error) error {
	return fmt.Errorf("blockwatch: cannot decode column '%s' into struct field of type %s: %v",
		name, val.Type().String(), err)
}

func makeColumnError(name string, col, row int, err error) error {
	return fmt.Errorf("blockwatch: cannot decode column '%s' [%d:%d]: %v", name, col, row, err)
}

func makeColumnMissingError(name string, col, row int) error {
	return fmt.Errorf("blockwatch: missing column '%s' [%d:%d]", name, col, row)
}
