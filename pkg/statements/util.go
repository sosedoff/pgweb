package statements

import (
	"fmt"
	"strings"
)

func JoinRecord(record []string) string {
	return fmt.Sprintf(`(%s)`, strings.Join(record, `, `))
}

func CreateBinding(fieldsNum int, rowNum int) string {
	rowBindings := make([]string, rowNum)
	for i := 0; i < rowNum; i++ {
		rowBindings[i] = JoinRecord(strings.Split(strings.Repeat(`?`, fieldsNum), ``))
	}
	return strings.Join(rowBindings, `, `)
}

func GenerateBulkInsertQuery(table string, header []string, rowNum int) string {
	return fmt.Sprintf(`INSERT INTO %s %s VALUES %s;`, table, JoinRecord(header), CreateBinding(len(header), rowNum))
}

func Flatten(records [][]string) []interface{} {
	rowNum := len(records)
	fieldsNum := len(records[0])
	flattenList := make([]interface{}, rowNum*fieldsNum)
	for i := 0; i < rowNum; i++ {
		for j := 0; j < fieldsNum; j++ {
			flattenList[i*fieldsNum+j] = records[i][j]
		}
	}
	return flattenList
}

func CreateNewTableQuery(table string, header []string) string {
	newHeader := make([]string, len(header)+1)
	for i, field := range header {
		newHeader[i] = fmt.Sprintf(`%s TEXT`, field)
	}
	newHeader[len(header)] = `id SERIAL PRIMARY KEY`
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s %s;`, table, JoinRecord(newHeader))
}
