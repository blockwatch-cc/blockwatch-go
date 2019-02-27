// Copyright (c) 2019 KIDTSUNAMI
// Author: alex@kidtsunami.com
//
package blockwatch

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type DatasetListParams struct {
	Limit  int
	Cursor string
}

func (p DatasetListParams) Query() url.Values {
	q := url.Values{}
	if p.Cursor != "" {
		q.Add("cursor", p.Cursor)
	}
	if p.Limit > 0 {
		q.Add("limit", strconv.Itoa(p.Limit))
	}
	return q
}

func (p DatasetListParams) Url(dbcode string) string {
	return fmt.Sprintf("databases/%s/codes.json?%s", dbcode, p.Query().Encode())
}

type Dataset struct {
	Database      string      `json:"database_code"`
	Dataset       string      `json:"dataset_code"`
	Type          string      `json:"type"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	Columns       []Datafield `json:"columns"`
	FilterFields  []string    `json:"filters"`
	PrimaryFields []string    `json:"primary_key"`
}

func (c *Client) ListDatasets(ctx context.Context, dbcode string, params DatasetListParams) ([]Dataset, error) {
	v := make([]Dataset, 0)
	err := c.Get(ctx, params.Url(dbcode), nil, &v)
	return v, err
}

func (c *Client) GetDataset(ctx context.Context, dbcode, setcode string) (*Dataset, error) {
	u := fmt.Sprintf("databases/%s/%s/metadata.json", dbcode, setcode)
	v := &Dataset{}
	err := c.Get(ctx, u, nil, v)
	return v, err
}
