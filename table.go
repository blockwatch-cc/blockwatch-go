// Copyright (c) 2018 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package blockwatch

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Table struct {
	Dataframe
	Limit  int    `json:"limit"`
	Count  int    `json:"count"`
	Cursor string `json:"cursor"`
	Error  *Error `json:"error"`
}

type TableParams struct {
	Columns []string
	Cursor  string
	Limit   int
	Filter  []*Filter
	Format  string
}

func (p TableParams) Query() url.Values {
	q := url.Values{}
	if len(p.Columns) > 0 && p.Columns[0] != "" {
		q.Add("columns", strings.Join(p.Columns, ","))
	}
	for _, v := range p.Filter {
		if v == nil {
			continue
		}
		v.AppendQuery(q)
	}
	if p.Cursor != "" {
		q.Add("cursor", p.Cursor)
	}
	if p.Limit > 0 {
		q.Add("limit", strconv.Itoa(p.Limit))
	}
	return q
}

func (p TableParams) Url(db, set string) string {
	if p.Format == "" {
		p.Format = "json"
	}
	return fmt.Sprintf("tables/%s/%s.%s?%s",
		db,
		set,
		p.Format,
		p.Query().Encode(),
	)
}

func (c *Client) GetTable(ctx context.Context, dbcode, setcode string, params TableParams) (*Table, error) {
	v := &Table{}
	err := c.Get(ctx, params.Url(dbcode, setcode), nil, v)
	if err != nil {
		return nil, err
	}
	// process streaming error
	if v.Error != nil {
		return v, v.Error
	}
	return v, nil
}

// TODO
// func (c *Client) StreamTable(ctx context.Context, params TableParams, fn func(r Row) error) error {
//  // slower, but from stdlib
//  // https://stackoverflow.com/questions/31794355/stream-large-json
//  // fetch chunks, stream-decode them and loop until limit is reached
// 	return nil
// }
