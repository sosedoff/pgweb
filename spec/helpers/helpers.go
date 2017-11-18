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
	CurrentDbSelector = "#current_database"
	ConnectionErrorSelector = "#connection_error"
)

func FillConnectionForm(page *agouti.Page, data map[string]string) {
	for selector, value := range data {
		page.Find(selector).Fill(value)
	}
}



func Screenshot(page *agouti.Page, name string) {
	page.Screenshot(fmt.Sprintf("_output/%s.png", name))
}

