package helpers

import (
	"fmt"

	"github.com/sclevine/agouti"
)

const (
	PgConnUrlSelector = "#connection_url"

	PgUserSelector = "#pg_user"
	PgPassSelector = "#pg_password"
	PgHostSelector = "#pg_host"
	PgPortSelector = "#pg_port"
	PgDbSelector   = "#pg_db"
	PgSslSelector  = "#connection_ssl"
)

const (
	CurrentDbSelector       = "#current_database"
	ConnectionErrorSelector = "#connection_error"
)

const (
	TableBookSelector = ".schema-table[data-id='public.books']"
)

const (
	TabRowsSelector       = "#table_content"
	TabStructureSelector  = "#table_stucture"
	TabIndexesSelector    = "#table_indexes"
	TabConstaintsSelector = "#table_constraints"
	TabQuerySelector      = "#table_query"
	TabHistorySelector    = "#table_history"
	TabActivitySelector   = "#table_activity"
	TabConnectionSelector = "#table_connection"
)

func FillConnectionForm(page *agouti.Page, data map[string]string) {
	for selector, value := range data {
		page.Find(selector).Fill(value)
	}
}

func ConnectByStandardTab(page *agouti.Page) {
	initVars()

	data := map[string]string{
		PgHostSelector: ServerHost,
		PgPortSelector: ServerPort,
		PgUserSelector: ServerUser,
		PgPassSelector: ServerPassword,
		PgDbSelector:   ServerDatabase,
		PgSslSelector:  "disable",
	}

	FillConnectionForm(page, data)
	page.FindByButton("Connect").Click()
}

func Screenshot(page *agouti.Page, name string) {
	page.Screenshot(fmt.Sprintf("_output/%s.png", name))
}
