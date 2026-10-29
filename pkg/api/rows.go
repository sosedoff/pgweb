package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/command"
)

// rowRequest is the payload for row insert, update and delete operations.
type rowRequest struct {
	// Key holds the primary key values identifying the target row.
	Key map[string]interface{} `json:"key"`

	// Values holds the column values to write.
	Values map[string]interface{} `json:"values"`
}

// rowEditingEnabled returns true when row modifications are available.
func rowEditingEnabled() bool {
	return !command.Opts.DisableRowEditing && !command.Opts.ReadOnly
}

// requireRowEditing blocks row modification endpoints when the feature is off.
func requireRowEditing() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rowEditingEnabled() {
			errorResponse(c, http.StatusForbidden, errRowEditingDisabled)
			return
		}

		c.Next()
	}
}

// parseRowRequest reads and validates the JSON body of a row request.
func parseRowRequest(c *gin.Context) (*rowRequest, error) {
	req := rowRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, errInvalidRowPayload
	}

	return &req, nil
}

// parseFilters decodes the JSON encoded list of structured row filters.
func parseFilters(value string) ([]client.Filter, error) {
	if value == "" {
		return nil, nil
	}

	filters := []client.Filter{}
	if err := json.Unmarshal([]byte(value), &filters); err != nil {
		return nil, errInvalidFilters
	}

	return filters, nil
}

// GetTableMetadata renders the column metadata used to build row forms and filters
func GetTableMetadata(c *gin.Context) {
	meta, err := DB(c).TableMetadata(c.Params.ByName("table"))
	if err != nil {
		badRequest(c, err)
		return
	}

	// Row editing availability is a server-side decision, so it is reported
	// alongside the metadata to keep the frontend logic simple.
	if !rowEditingEnabled() {
		meta = disableWrites(meta)
	}

	successResponse(c, meta)
}

// CreateTableRow inserts a new row into the table
func CreateTableRow(c *gin.Context) {
	req, err := parseRowRequest(c)
	if err != nil {
		badRequest(c, err)
		return
	}

	res, err := DB(c).InsertRow(c.Params.ByName("table"), req.Values)
	if err != nil {
		badRequest(c, err)
		return
	}

	successResponse(c, rowResponse(res))
}

// UpdateTableRow updates a single existing row
func UpdateTableRow(c *gin.Context) {
	req, err := parseRowRequest(c)
	if err != nil {
		badRequest(c, err)
		return
	}

	res, err := DB(c).UpdateRow(c.Params.ByName("table"), req.Key, req.Values)
	if err != nil {
		badRequest(c, err)
		return
	}

	successResponse(c, rowResponse(res))
}

// DeleteTableRow removes a single existing row
func DeleteTableRow(c *gin.Context) {
	req, err := parseRowRequest(c)
	if err != nil {
		badRequest(c, err)
		return
	}

	res, err := DB(c).DeleteRow(c.Params.ByName("table"), req.Key)
	if err != nil {
		badRequest(c, err)
		return
	}

	successResponse(c, rowResponse(res))
}

// rowResponse formats the affected row so that the frontend can patch the grid
// in place instead of reloading the whole table.
func rowResponse(res *client.Result) gin.H {
	response := gin.H{"columns": res.Columns, "rows": res.Rows}

	if len(res.Rows) > 0 {
		response["row"] = res.Format()[0]
	}

	return response
}

// disableWrites returns a copy of the metadata with write support turned off.
func disableWrites(meta *client.TableMetadata) *client.TableMetadata {
	copied := *meta
	copied.Insertable = false
	copied.Editable = false
	return &copied
}
