## blockwatch-go – Official Go SDK for the Blockwatch Data API

The official [Blockwatch](https://blockwatch.cc) Go client library.

To use this SDK you need a free API key and a free or payed database subscription.

### Installation

```sh
go get -u blockwatch.cc/blockwatch-go
```

Then import, using

```go
import (
	"blockwatch.cc/blockwatch-go"
)
```

### Documentation

For a comprehensive coverage of all API features read the [Blockwatch Data API](https://blockwatch.cc/docs/api) documentation.

Below are a few examples on how to use the Go SDK to access the Blockwatch Data API.

### Authentication

Authentication for the Blockwatch Data API works by using API keys as secret tokens. You can get your API key by signing up for a [free Blockwatch account](https://blockwatch.cc/account/signup) and then creating your personal API key on your [account settings page](https://blockwatch.cc/account/profile#apikey).

You also need to have an active subscription to the databases you like to query. You will get `404 Not Found` errors if you try to access a databases without subscription.


### Initializing the Go SDK Client

For convenient access to the Blockwatch Data API, the Go SDK defines a `Client` that exports all relevant functions. To create a new client object with default configuration call:

```go
c, err := blockwatch.NewClient("MY_API_KEY", nil)
```

The default configuration should work just fine, but if you need special timeouts, proxy or TLS settings you may specify a separate `ConnConfig` struct.

```go
type ConnConfig struct {
	// HTTP tuning parameters
	DialTimeout           time.Duration
	KeepAlive             time.Duration
	IdleConnTimeout       time.Duration
	ResponseHeaderTimeout time.Duration
	ExpectContinueTimeout time.Duration
	MaxIdleConns          int

	// Proxy specifies to connect through a SOCKS 5 proxy server.  It may
	// be an empty string if a proxy is not required.
	Proxy string

	// ProxyUser is an optional username to use for the proxy server if it
	// requires authentication.  It has no effect if the Proxy parameter
	// is not set.
	ProxyUser string

	// ProxyPass is an optional password to use for the proxy server if it
	// requires authentication.  It has no effect if the Proxy parameter
	// is not set.
	ProxyPass string

	// TLS configuration options
	ServerName         string
	AllowInsecureCerts bool
	TLSMinVersion      int
	TLSMaxVersion      int
	RootCaCerts        []string
	RootCaCertsFile    string
	ClientCert         []string
	ClientCertFile     string
	ClientKey          []string
	ClientKeyFile      string
}
```

### Listing Databases

To fetch a list of all databases you're currently subscribed to, call:

```go
dbs, err := c.ListDatabases(ctx, blockwatch.DatabaseListParams{})
```

### Getting Dataset Information

To get the list of all datasets in a particular database, call

```go
sets, err := c.ListDatasets(ctx, "BTC", blockwatch.DatasetListParams{})
```

While the above call returns just the code and name of all datasets, you can fetch full details including the list of data fields with

```go
set, err := c.GetDataset(ctx,  "BTC", "BLOCK")
```

### Fetching Dataset Contents

The Blockwatch Data API supports two different kinds of datasets, **tables** and **time-series** with slightly different query arguments and semantics. Both contain a matrix of columns and rows which you can access using the same kind of functions.

Time series may contain at most one row of entries per unique timestamp. They have a default sampling frequency, are naturally ordered by timestamp and you may only filter them by time in most cases. You can request a different frequency using the `collapse` parameter in which case rows will be automatically aggregated by time.

Tables can contain arbitrary data and are are ordered by a primary key, in most cases a unique row id. Usually you can filter tables by most of their columns. Please consult the dataset spec for details on which data fields are filterable.

The SDK keeps the raw data response and lazy-unmarshals data rows or columns only when accessed.


#### Getting Data from Tables

```go
# using default parameters (limit = 500, all columns)
table, err := c.GetTable(ctx, "BTC", "BLOCK", blockwatch.TableParams{})

# with limit and filters
table, err = c.GetTable(ctx, "BTC", "BLOCK", blockwatch.TableParams{
	Limit:   10,
	Columns: "time,height,hash,volume",
	Filter: []*blockwatch.Filter{
		&blockwatch.Filter{
			Field: "height",
			Mode: blockwatch.FilterModeGte,
			Value: "500000",
		},
		blockwatch.NewFilter("volume", blockwatch.FilterModeRange, "10,100"),
	},
})
```

#### Getting Time-series Data

```go
# using default parameters (limit = 500, all columns)
series, err := c.GetSeries(ctx, "BITFINEX:OHLCV", "BTC_USD", blockwatch.SeriesParams{})

# get hourly data from the last 24 hours
series, err = c.GetSeries(ctx, "BITFINEX:OHLCV", "BTC_USD", blockwatch.SeriesParams{
	Limit:     24,
	Columns:   "time,open,close,vwap,vol_base",
	Collapse:  blockwatch.CollapseModeOneHour,
	Order:     OrderModeDesc,
	StartDate: time.Now().Add(-24*time.Hour),
	EndDate:   time.Now(),
})
```

### Decoding Data into Structs

To process data row-by-row it is often convenient to extract each row into a Go struct first. We've defined a couple of common structs for Blockwatch market and blockchain databases, but you can also define your own structs. This makes sense if you like to limit the number of struct fields for memory efficiency.

For extraction the SDK uses Go's built-in type reflection and matches **column codes** to **struct field tags**. Per default the SDK looks for `json` struct tags, but you can change that by setting `blockwatch.TagName = "mytag"`.

```go
// assuming you have fetched data from BTC/BLOCK into table
err := table.ForEach(func(r blockwatch.Row) error {
	var block blockwatch.Block

	// decode row into struct fields (uses json struct tags to match column codes)
	if err := r.Decode(&block); err != nil {
		return err
	}

	// handle block data here

	return nil
})
```

### Decoding Columns as Slices

Sometimes it's more efficient to process data in column vectors. To support this mode you may extract an entire column in one step:

```go
// assuming you have fetched data from BTC/BLOCK into table, you may now
// decode each column as slice (Note that an interface to slice is returned)
idx, col, err := table.Column("hash")

// to work with a native slice you'll have to cast the interface to type []<T>;
// either choose a static (i.e. implementation time) approach or a dynamic one

// static (requires you know the type for column 'hash' in advance)
slice, _ := col.([][]byte)

// or dynamic (using the actual column type in a switch statement)
switch table.Columns[idx].Type {
	case blockwatch.FieldTypeString:
		slice, _ := col.([]string)

	case blockwatch.FieldTypeBytes:
		slice, _ := col.([][]byte)

	case blockwatch.FieldTypeDate, blockwatch.FieldTypeDatetime:
		slice, _ := col.([]time.Time)

	case blockwatch.FieldTypeBoolean:
		slice, _ := col.([]bool)

	case blockwatch.FieldTypeFloat64:
		slice, _ := col.([]float64)

	case blockwatch.FieldTypeInt64:
		slice, _ := col.([]int64)

	case blockwatch.FieldTypeUint64:
		slice, _ := col.([]uint64)

	default:
		// handle unsupported column type
}
```

### Cursoring through large result sets

A query can match millions of rows, but for efficiency reasons we limit each result to at most 50,000 rows. A result contains a `cursor` value that allows you to fetch the next chunk of rows right after the current one in a subsequent query. When a result contains no more data you know that you've reached the end of a table. Because most tables grow in real-time you can also store the latest cursor and poll for new data after a while.

```go
params := blockwatch.TableParams{}

for {
	table, err := blockwatch.GetTable(ctx, "BTC", "BLOCK", params)
	// handle error if necessary

	// handle data here

	// prepare for next iteration
	params.Cursor = table.Cursor
}

```


### Gracefully handling rate-limits

To avoid excessive overload of our API we limit the rate at which we process your requests. This means your program may from time to time run into a rate limit. To let you gracefully handle retries by waiting until a rate limit resets, we expose the deadline and a done channel much like Go's network context does. Here's how you may use this feature:

```go
var (
	table *blockwatch.Table
	err   error
)
for {
	table, err = blockwatch.GetTable(ctx, "BTC", "BLOCK", blockwatch.TableParams{})
	if err != nil {
		if e, ok := blockwatch.IsRateLimited(err); ok {
			fmt.Printf("Rate limited, waiting for %s\n", e.Deadline())
			select {
			case <-ctx.Done():
				// wait until external context is canceled
				err = ctx.Err()
			case <-e.Done():
				// wait until rate limit reset and retry
				continue
			}
		}
	}
	break
}

// handle error and/or result here

```

## License

The MIT License (MIT) Copyright (c) 2020 Blockwatch Data Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is furnished
to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.