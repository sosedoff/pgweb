package client

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
)

type Row []interface{}

type Result struct {
	Columns []string `json:"columns"`
	Rows    []Row    `json:"rows"`
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
	records := []map[string]interface{}{}

	for _, row := range res.Rows {
		record := map[string]interface{}{}
		for i, col := range res.Columns {
			record[col] = row[i]
		}
		records = append(records, record)
	}

	data, _ := json.Marshal(records)
	return data
}
