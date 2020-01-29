// Copyright (c) 2020 Blockwatch Data Inc.
// Author: alex@blockwatch.cc
//
package blockwatch

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type DatabaseListParams struct {
	Limit  int
	Cursor string
}

func (p DatabaseListParams) Query() url.Values {
	q := url.Values{}
	if p.Cursor != "" {
		q.Add("cursor", p.Cursor)
	}
	if p.Limit > 0 {
		q.Add("limit", strconv.Itoa(p.Limit))
	}
	return q
}

func (p DatabaseListParams) Url() string {
	return fmt.Sprintf("databases?%s", p.Query().Encode())
}

type Database struct {
	Id                string    `json:"database_id"`
	AuthorId          string    `json:"author_id"`
	Code              string    `json:"code"`
	Name              string    `json:"name"`
	DatasetType       string    `json:"type"`
	State             string    `json:"state"`
	Description       string    `json:"description"`
	Documentation     string    `json:"documentation"`
	ImageId           string    `json:"imageId"`
	IsPremium         bool      `json:"is_premium"`
	HasSample         bool      `json:"has_sample"`
	DeliveryFrequency string    `json:"delivery_frequency"`
	DataFrequency     string    `json:"data_frequency"`
	ReportingLag      string    `json:"reporting_lag"`
	History           string    `json:"history"`
	Coverage          string    `json:"coverage"`
	Labels            string    `json:"labels"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Subscribed        bool      `json:"subscribed"`
}

type DatabaseList struct {
	Meta struct {
		Count  int    `json:"count"`
		Cursor string `json:"cursor"`
	} `json:"meta"`
	Databases []*Database `json:"databases"`
}

func (c *Client) ListDatabases(ctx context.Context, params DatabaseListParams) (*DatabaseList, error) {
	v := &DatabaseList{}
	err := c.Get(ctx, params.Url(), nil, v)
	return v, err
}
