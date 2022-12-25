package client

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/sosedoff/pgweb/pkg/command"
)

const (
	ObjTypeTable            = "table"
	ObjTypeView             = "view"
	ObjTypeMaterializedView = "materialized_view"
	ObjTypeSequence         = "sequence"
	ObjTypeFunction         = "function"
)

type (
	// Row represents a single row of data
	Row []interface{}

	// RowsOptions contains a list of parameters for table browsing requests
	RowsOptions struct {
		Where      string // Custom filter
		Offset     int    // Number of rows to skip
		Limit      int    // Number of rows to fetch
		SortColumn string // Column to sort by
		SortOrder  string // Sort direction (ASC, DESC)
	}

	Pagination struct {
		Rows    int64 `json:"rows_count"`
		Page    int64 `json:"page"`
		Pages   int64 `json:"pages_count"`
		PerPage int64 `json:"per_page"`
	}

	Result struct {
		Pagination *Pagination  `json:"pagination,omitempty"`
		Columns    []string     `json:"columns"`
		Rows       []Row        `json:"rows"`
		Stats      *ResultStats `json:"stats,omitempty"`
	}

	ResultStats struct {
		ColumnsCount    int       `json:"columns_count"`
		RowsCount       int       `json:"rows_count"`
		RowsAffected    int64     `json:"rows_affected"`
		QueryStartTime  time.Time `json:"query_start_time"`
		QueryFinishTime time.Time `json:"query_finish_time"`
		QueryDuration   int64     `json:"query_duration_ms"`
	}

	Object struct {
		OID  string `json:"oid"`
		Name string `json:"name"`
	}

	Objects struct {
		Tables            []Object `json:"table"`
		Views             []Object `json:"view"`
		MaterializedViews []Object `json:"materialized_view"`
		Functions         []Object `json:"function"`
		Sequences         []Object `json:"sequence"`
	}
)

// Due to big int number limitations in javascript, numbers should be encoded
// as strings so they could be properly loaded on the frontend.
func (res *Result) PostProcess() {
	for i, row := range res.Rows {
		for j, col := range row {
			if col == nil {
				continue
			}

			switch val := col.(type) {
			case int64:
				if val < -9007199254740991 || val > 9007199254740991 {
					res.Rows[i][j] = strconv.FormatInt(col.(int64), 10)
				}
			case float64:
				// json.Marshal panics when dealing with NaN/Inf values
				// issue: https://github.com/golang/go/issues/25721
				if math.IsNaN(val) {
					res.Rows[i][j] = nil
					break
				}

				if val < -999999999999999 || val > 999999999999999 {
					res.Rows[i][j] = strconv.FormatFloat(val, 'e', -1, 64)
				}
			case string:
				if hasBinary(val, 8) && BinaryCodec != CodecNone {
					res.Rows[i][j] = encodeBinaryData([]byte(val), BinaryCodec)
				}
			case time.Time:
				// RFC 3339 is clear that years are 4 digits exactly.
				// See golang.org/issue/4556#c15 for more discussion.
				if val.Year() < 0 || val.Year() >= 10000 {
					res.Rows[i][j] = "ERR: INVALID_DATE"
				} else {
					res.Rows[i][j] = val
				}
			}
		}
	}
}

func (res *Result) Format() []map[string]interface{} {
	items := make([]map[string]interface{}, len(res.Rows))

	for rowIdx, row := range res.Rows {
		item := make(map[string]interface{})
		for i, c := range res.Columns {
			item[c] = row[i]
		}

		items[rowIdx] = item
	}

	return items
}

func (res *Result) CSV() []byte {
	buff := &bytes.Buffer{}
	writer := csv.NewWriter(buff)

	if err := writer.Write(res.Columns); err != nil {
		log.Printf("result csv write error: %v\n", err)
	}

	for _, row := range res.Rows {
		record := make([]string, len(res.Columns))

		for i, item := range row {
			switch v := item.(type) {
			case time.Time:
				record[i] = v.Format("2006-01-02 15:04:05")
			case nil:
				record[i] = ""
			default:
				record[i] = fmt.Sprintf("%v", item)
			}
		}

		err := writer.Write(record)
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	writer.Flush()
	return buff.Bytes()
}

func (res *Result) JSON() []byte {
	var data []byte

	if command.Opts.DisablePrettyJSON {
		data, _ = json.Marshal(res.Format())
	} else {
		data, _ = json.MarshalIndent(res.Format(), "", " ")
	}

	return data
}

func ObjectsFromResult(res *Result) map[string]*Objects {
	objects := map[string]*Objects{}

	for _, row := range res.Rows {
		oid := row[0].(string)
		schema := row[1].(string)
		name := row[2].(string)
		objectType := row[3].(string)

		if objects[schema] == nil {
			objects[schema] = &Objects{
				Tables:            []Object{},
				Views:             []Object{},
				MaterializedViews: []Object{},
				Functions:         []Object{},
				Sequences:         []Object{},
			}
		}

		obj := Object{OID: oid, Name: name}

		switch objectType {
		case ObjTypeTable:
			objects[schema].Tables = append(objects[schema].Tables, obj)
		case ObjTypeView:
			objects[schema].Views = append(objects[schema].Views, obj)
		case ObjTypeMaterializedView:
			objects[schema].MaterializedViews = append(objects[schema].MaterializedViews, obj)
		case ObjTypeFunction:
			objects[schema].Functions = append(objects[schema].Functions, obj)
		case ObjTypeSequence:
			objects[schema].Sequences = append(objects[schema].Sequences, obj)
		}
	}

	return objects
}
