// Copyright (c) 2019 KIDTSUNAMI
// Author: alex@kidtsunami.com

// Blockwatch Data API example

package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"blockwatch.cc/blockwatch-go"
)

var (
	flags    = flag.NewFlagSet("blockwatch", flag.ContinueOnError)
	verbose  bool
	key      string
	cmd      string
	code     string
	columns  string
	collapse string
	limit    int
)

var ()

var cmdinfo = `
Available Commands:
  list-dbs            list all databases
  list-sets           list all datasets in a given database
  list-fields         list all datafields in a given dataset
  series              show time-series data
  table               show table data
  parse-table         parse table into blockchain BLOCK struct
  parse-table-column  parse first table column
`

func init() {
	flags.Usage = func() {}
	flags.BoolVar(&verbose, "v", false, "be verbose")
	flags.StringVar(&key, "apikey", "", "Blockwatch API key")
	flags.IntVar(&limit, "limit", 5, "row limit")
	flags.StringVar(&columns, "columns", "", "list of columns")
	flags.StringVar(&collapse, "collapse", "", "collapse mode for time-series (1m, 1h, 1d)")
}

func printhelp() {
	fmt.Println("Usage:\n  blockwatch [flags] [command] [args]")
	fmt.Println(cmdinfo)
	fmt.Println("Flags:")
	flags.PrintDefaults()
	fmt.Println()
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := flags.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			printhelp()
			return nil
		}
		return err
	}

	switch flags.NArg() {
	case 0:
		return fmt.Errorf("Missing command.")
	case 1:
		cmd = flags.Arg(0)
	default:
		cmd = flags.Arg(0)
		code = flags.Arg(1)
	}

	c, err := blockwatch.NewClient(key, nil)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch cmd {
	case "list-dbs":
		return listDatabases(ctx, c)
	case "list-sets":
		return listDatasets(ctx, c)
	case "list-fields":
		return listDatafields(ctx, c)
	case "table":
		return fetchTable(ctx, c)
	case "series":
		return fetchSeries(ctx, c)
	case "parse-table":
		return parseTable(ctx, c)
	case "parse-table-column":
		return parseTableColumn(ctx, c)
	default:
		return fmt.Errorf("unsupported command %s", cmd)
	}
}

func listDatabases(ctx context.Context, c *blockwatch.Client) error {
	dbs, err := c.ListDatabases(ctx, blockwatch.DatabaseListParams{
		Limit: limit,
	})
	if err != nil {
		return err
	}
	if dbs.Meta.Count == 0 {
		fmt.Println("No databases found. Are you subscribed?")
		return nil
	}
	fmtstr := "%-3v %-20s %-45s %-10s\n"
	fmt.Printf(fmtstr, "#", "Code", "Name", "Type")
	for i, v := range dbs.Databases {
		fmt.Printf(fmtstr, i+1, v.Code, v.Name, v.DatasetType)
	}
	return nil
}

func listDatasets(ctx context.Context, c *blockwatch.Client) error {
	if code == "" {
		return fmt.Errorf("missing database code")
	}
	sets, err := c.ListDatasets(ctx, code, blockwatch.DatasetListParams{
		Limit: limit,
	})
	if err != nil {
		return err
	}
	if len(sets) == 0 {
		fmt.Println("No dataset found. Are you subscribed?")
		return nil
	}
	fmtstr := "%-3v %-30s %-40s\n"
	fmt.Printf(fmtstr, "#", "Code", "Name")
	for i, v := range sets {
		fmt.Printf(fmtstr, i+1, v.Database+"/"+v.Dataset, v.Name)
	}
	return nil
}

func listDatafields(ctx context.Context, c *blockwatch.Client) error {
	cf := strings.Split(code, "/")
	if len(cf) != 2 {
		return fmt.Errorf("invalid dataset code")
	}
	set, err := c.GetDataset(ctx, cf[0], cf[1])
	if err != nil {
		return err
	}
	fmtstr := "%-3v %-20s %-25s %-10s %4s %9s\n"
	fmt.Printf(fmtstr, "#", "Code", "Name", "Type", "Primary", "Filterable")
	for i, v := range set.Columns {
		var isprimary, isfilter string
		if contains(set.PrimaryFields, v.Code) {
			isprimary = "*"
		}
		if contains(set.FilterFields, v.Code) {
			isfilter = "*"
		}
		fmt.Printf(fmtstr, i+1, v.Code, v.Name, v.Type, isprimary, isfilter)
	}
	return nil
}

func fetchTable(ctx context.Context, c *blockwatch.Client) error {
	cf := strings.Split(code, "/")
	if len(cf) != 2 {
		return fmt.Errorf("invalid dataset code")
	}
	table, err := c.GetTable(ctx, cf[0], cf[1], blockwatch.TableParams{
		Limit:   limit,
		Columns: strings.Split(columns, ","),
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s", dumpData(table.Dataframe))
	return nil
}

func parseTable(ctx context.Context, c *blockwatch.Client) error {
	cf := strings.Split(code, "/")
	if len(cf) != 2 {
		return fmt.Errorf("invalid dataset code")
	}
	table, err := c.GetTable(ctx, cf[0], cf[1], blockwatch.TableParams{
		Limit:   limit,
		Columns: strings.Split(columns, ","),
	})
	if err != nil {
		return err
	}
	var block blockwatch.Block
	err = table.ForEach(func(r blockwatch.Row) error {
		if err := r.Decode(&block); err != nil {
			return err
		}
		fmt.Printf("Decoded block %d %s\n", block.Height, block.Timestamp)
		return nil
	})
	return err
}

func parseTableColumn(ctx context.Context, c *blockwatch.Client) error {
	cf := strings.Split(code, "/")
	if len(cf) != 2 {
		return fmt.Errorf("invalid dataset code")
	}
	table, err := c.GetTable(ctx, cf[0], cf[1], blockwatch.TableParams{
		Limit:   limit,
		Columns: strings.Split(columns, ","),
	})
	if err != nil {
		return err
	}

	// decode column as slice (will return an interface to slice)
	_, col, err := table.Column(table.Columns[0].Code)
	if err != nil {
		return err
	}

	// to work with the slice natively you'll have to cast the interface
	// to slice of type table.Columns[0].Type (here's a fully dynamic approach,
	// in case you know at implementation time you can make it static)
	switch typ := table.Columns[0].Type; typ {
	case blockwatch.FieldTypeString:
		slice, _ := col.([]string)
		fmt.Printf("Decoded %s column '%s' with %d values, first is %v, last is %v\n",
			typ, table.Columns[0].Code, len(slice), slice[0], slice[len(slice)-1])

	case blockwatch.FieldTypeBytes:
		slice, _ := col.([][]byte)
		fmt.Printf("Decoded %s column '%s' with %d values, first is %x, last is %x\n",
			typ, table.Columns[0].Code, len(slice), slice[0], slice[len(slice)-1])

	case blockwatch.FieldTypeDate, blockwatch.FieldTypeDatetime:
		slice, _ := col.([]time.Time)
		fmt.Printf("Decoded %s column '%s' with %d values, first is %s, last is %s\n",
			typ, table.Columns[0].Code, len(slice), slice[0], slice[len(slice)-1])

	case blockwatch.FieldTypeBoolean:
		slice, _ := col.([]bool)
		fmt.Printf("Decoded %s column '%s' with %d values, first is %v, last is %v\n",
			typ, table.Columns[0].Code, len(slice), slice[0], slice[len(slice)-1])

	case blockwatch.FieldTypeFloat64:
		slice, _ := col.([]float64)
		fmt.Printf("Decoded %s column '%s' with %d values, first is %v, last is %v\n",
			typ, table.Columns[0].Code, len(slice), slice[0], slice[len(slice)-1])

	case blockwatch.FieldTypeInt64:
		slice, _ := col.([]int64)
		fmt.Printf("Decoded %s column '%s' with %d values, first is %v, last is %v\n",
			typ, table.Columns[0].Code, len(slice), slice[0], slice[len(slice)-1])

	case blockwatch.FieldTypeUint64:
		slice, _ := col.([]uint64)
		fmt.Printf("Decoded %s column '%s' with %d values, first is %v, last is %v\n",
			typ, table.Columns[0].Code, len(slice), slice[0], slice[len(slice)-1])

	default:
		return fmt.Errorf("Unsupported column type %s", typ)
	}

	return nil
}

func fetchSeries(ctx context.Context, c *blockwatch.Client) error {
	cf := strings.Split(code, "/")
	if len(cf) != 2 {
		return fmt.Errorf("invalid dataset code")
	}
	series, err := c.GetSeries(ctx, cf[0], cf[1], blockwatch.SeriesParams{
		Limit:    limit,
		Columns:  strings.Split(columns, ","),
		Collapse: blockwatch.ParseCollapseModeIgnoreError(collapse),
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s", dumpData(series.Dataframe))
	return nil
}

// helper functions

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func dumpData(t blockwatch.Dataframe) string {
	cols, rows := len(t.Columns), len(t.Data)
	dump := make([][]string, cols)
	sz := make([]int, cols)
	for j := 0; j < cols; j++ {
		dump[j] = make([]string, rows)
		sz[j] = len(t.Columns[j].Name)
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			val, err := t.FieldAt(j, i)
			if err != nil {
				fmt.Printf("Field %d/%d: %v\n", j, i, err)
			}
			dump[j][i] = toString(val)
			sz[j] = max(sz[j], len(dump[j][i]))
		}
	}
	row := make([]string, cols)
	for j := 0; j < cols; j++ {
		row[j] = fmt.Sprintf("%[2]*[1]s", t.Columns[j].Name, -sz[j])
	}
	grid := make([]string, rows+2)
	grid[0] = "| " + strings.Join(row, " | ") + " |"
	for j := 0; j < cols; j++ {
		row[j] = strings.Repeat("-", sz[j])
	}
	grid[1] = "|-" + strings.Join(row, "-|-") + "-|"
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			row[j] = fmt.Sprintf("%[2]*[1]s", dump[j][i], -sz[j])
		}
		grid[i+2] = "| " + strings.Join(row, " | ") + " |"
	}
	return strings.Join(grid, "\n")
}

var stringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

func toString(t interface{}) string {
	val := reflect.Indirect(reflect.ValueOf(t))
	if !val.IsValid() {
		return ""
	}
	if val.Type().Implements(stringerType) {
		return t.(fmt.Stringer).String()
	}
	if s, err := toRawString(val.Interface()); err == nil {
		return s
	}
	return fmt.Sprintf("%v", val.Interface())
}

func toRawString(t interface{}) (string, error) {
	val := reflect.Indirect(reflect.ValueOf(t))
	if !val.IsValid() {
		return "", nil
	}
	typ := val.Type()
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(val.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, val.Type().Bits()), nil
	case reflect.String:
		return val.String(), nil
	case reflect.Bool:
		return strconv.FormatBool(val.Bool()), nil
	case reflect.Array:
		if typ.Elem().Kind() != reflect.Uint8 {
			break
		}
		// [...]byte
		var b []byte
		if val.CanAddr() {
			b = val.Slice(0, val.Len()).Bytes()
		} else {
			b = make([]byte, val.Len())
			reflect.Copy(reflect.ValueOf(b), val)
		}
		return hex.EncodeToString(b), nil
	case reflect.Slice:
		if typ.Elem().Kind() != reflect.Uint8 {
			break
		}
		// []byte
		b := val.Bytes()
		return hex.EncodeToString(b), nil
	}
	return "", fmt.Errorf("no method for converting type %s (%v) to string", typ.String(), val.Kind())
}

func max(x, y int) int {
	if x < y {
		return y
	} else {
		return x
	}
}
