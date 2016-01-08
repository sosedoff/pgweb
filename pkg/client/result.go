package client

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type Row []interface{}

type Result struct {
	Columns []string `json:"columns"`
	Rows    []Row    `json:"rows"`
}

// Due to big int number limitations in javascript, numbers should be encoded
// as strings so they could be properly loaded on the frontend.
func (res *Result) PrepareBigints() {
	for i, row := range res.Rows {
		for j, col := range row {
			if col == nil {
				continue
			}

			switch reflect.TypeOf(col).Kind() {
			case reflect.Int64:
				val := col.(int64)
				if val < -9007199254740991 || val > 9007199254740991 {
					res.Rows[i][j] = strconv.FormatInt(col.(int64), 10)
				}
			case reflect.Float64:
				val := col.(float64)
				if val < -999999999999999 || val > 999999999999999 {
					res.Rows[i][j] = strconv.FormatFloat(val, 'e', -1, 64)
				}
			}
		}
	}
}

func (res *Result) Format() []map[string]interface{} {
	var items []map[string]interface{}

	for _, row := range res.Rows {
		item := make(map[string]interface{})

		for i, c := range res.Columns {
			item[c] = row[i]
		}

		items = append(items, item)
	}

	return items
}

func (res *Result) CSV() []byte {
	buff := &bytes.Buffer{}
	writer := csv.NewWriter(buff)

	writer.Write(res.Columns)

	for _, row := range res.Rows {
		record := make([]string, len(res.Columns))

		for i, item := range row {
			if item != nil {
				record[i] = fmt.Sprintf("%v", item)
			} else {
				record[i] = ""
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
	data, _ := json.Marshal(res.Format())
	return data
}
