// Copyright (c) 2019 KIDTSUNAMI
// Author: alex@kidtsunami.com

package blockwatch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Series struct {
	Dataframe
	Collapse  CollapseMode `json:"collapse"`
	Order     OrderMode    `json:"order"`
	StartDate time.Time    `json:"start_date"`
	EndDate   time.Time    `json:"end_date"`
	Limit     int          `json:"limit"`
	Count     int          `json:"count"`
	Error     *Error       `json:"error"`
}

// convert timestamps on unmarshal
func (s *Series) UnmarshalJSON(data []byte) error {
	series := struct {
		Dataframe
		Collapse CollapseMode `json:"collapse"`
		Order    OrderMode    `json:"order"`
		Start    int64        `json:"start_date"`
		End      int64        `json:"end_date"`
		Limit    int          `json:"limit"`
		Count    int          `json:"count"`
		Error    *Error       `json:"error"`
	}{}
	if err := json.Unmarshal(data, &series); err != nil {
		return err
	}
	s.Dataframe = series.Dataframe
	s.Collapse = series.Collapse
	s.Order = series.Order
	s.StartDate = time.Unix(0, series.Start*1000000).UTC()
	s.EndDate = time.Unix(0, series.End*1000000).UTC()
	s.Limit = series.Limit
	s.Count = series.Count
	s.Error = series.Error
	return nil
}

type SeriesParams struct {
	Columns   []string
	Collapse  CollapseMode
	Order     OrderMode
	StartDate time.Time
	EndDate   time.Time
	Limit     int
	Format    string
	Filter    []*Filter
}

func (p SeriesParams) Query() url.Values {
	q := url.Values{}
	if len(p.Columns) > 0 && p.Columns[0] != "" {
		q.Add("columns", strings.Join(p.Columns, ","))
	}
	if p.Collapse.IsValid() {
		q.Add("collapse", p.Collapse.String())
	}
	if p.Order.IsValid() {
		q.Add("order", p.Order.String())
	}
	if !p.StartDate.IsZero() {
		q.Add("start_date", p.StartDate.Format(time.RFC3339))
	}
	if !p.EndDate.IsZero() {
		q.Add("end_date", p.EndDate.Format(time.RFC3339))
	}
	for _, v := range p.Filter {
		if v == nil {
			continue
		}
		v.AppendQuery(q)
	}
	if p.Limit > 0 {
		q.Add("limit", strconv.Itoa(p.Limit))
	}
	return q
}

func (p SeriesParams) Url(db, set string) string {
	if p.Format == "" {
		p.Format = "json"
	}
	return fmt.Sprintf("series/%s/%s.%s?%s",
		db,
		set,
		p.Format,
		p.Query().Encode(),
	)
}

func (c *Client) GetSeries(ctx context.Context, dbcode, setcode string, params SeriesParams) (*Series, error) {
	v := &Series{}
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
