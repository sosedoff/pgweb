package client

import (
	"fmt"
	"strings"
)

// filterOperator describes a single supported WHERE condition.
type filterOperator struct {
	// sql is the operator rendered into the statement.
	sql string

	// args is the number of bind parameters the operator consumes:
	// 0 for IS NULL style checks, 1 for binary operators and -1 for lists.
	args int
}

// filterOperators is the whitelist of operators available to the row filters.
// Keys match the values used by the frontend filter dropdown.
var filterOperators = map[string]filterOperator{
	"equal":      {sql: "=", args: 1},
	"not_equal":  {sql: "!=", args: 1},
	"greater":    {sql: ">", args: 1},
	"greater_eq": {sql: ">=", args: 1},
	"less":       {sql: "<", args: 1},
	"less_eq":    {sql: "<=", args: 1},
	"like":       {sql: "LIKE", args: 1},
	"ilike":      {sql: "ILIKE", args: 1},
	"in":         {sql: "IN", args: -1},
	"null":       {sql: "IS NULL", args: 0},
	"not_null":   {sql: "IS NOT NULL", args: 0},
}

// Filter is a single column condition used to narrow down table rows.
type Filter struct {
	Column string `json:"column"`
	Op     string `json:"op"`
	Value  string `json:"value"`
}

// FilterOperators returns the list of supported filter operator names.
func FilterOperators() []string {
	names := make([]string, 0, len(filterOperators))
	for name := range filterOperators {
		names = append(names, name)
	}
	return names
}

// buildFilters renders a list of filters into a parameterized SQL condition.
//
// Column names are validated against the table metadata before they are
// quoted, and every user supplied value is passed as a bind parameter so that
// filter input can never alter the statement structure.
//
// The startIdx argument is the number of bind parameters already used by the
// caller, placeholders are numbered from startIdx+1.
func buildFilters(filters []Filter, meta *TableMetadata, startIdx int) (string, []interface{}, error) {
	if len(filters) == 0 {
		return "", nil, nil
	}

	conditions := make([]string, 0, len(filters))
	args := []interface{}{}
	placeholder := startIdx

	for _, filter := range filters {
		column := strings.TrimSpace(filter.Column)
		if column == "" {
			return "", nil, fmt.Errorf("filter column is required")
		}

		if meta != nil && meta.Column(column) == nil {
			return "", nil, fmt.Errorf("unknown filter column: %s", column)
		}

		op, ok := filterOperators[filter.Op]
		if !ok {
			return "", nil, fmt.Errorf("unsupported filter operator: %s", filter.Op)
		}

		switch op.args {
		case 0:
			conditions = append(conditions, fmt.Sprintf("%s %s", quoteIdent(column), op.sql))
		case 1:
			placeholder++
			conditions = append(conditions, fmt.Sprintf("%s %s $%d", quoteIdent(column), op.sql, placeholder))
			args = append(args, filter.Value)
		default:
			values := splitListValue(filter.Value)
			if len(values) == 0 {
				return "", nil, fmt.Errorf("filter on %s requires at least one value", column)
			}

			marks := make([]string, len(values))
			for i, value := range values {
				placeholder++
				marks[i] = fmt.Sprintf("$%d", placeholder)
				args = append(args, value)
			}

			conditions = append(conditions, fmt.Sprintf("%s %s (%s)", quoteIdent(column), op.sql, strings.Join(marks, ", ")))
		}
	}

	return strings.Join(conditions, " AND "), args, nil
}

// splitListValue splits a comma separated IN filter value into its elements.
func splitListValue(value string) []string {
	items := []string{}

	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		items = append(items, item)
	}

	return items
}
