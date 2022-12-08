package static

import (
	"embed"
	"net/http"
	"os"
)

//go:embed img/* js/* css/* fonts/*
//go:embed index.html
var assets embed.FS

func GetFilesystem() http.FileSystem {
	if os.Getenv("PGWEB_ASSETS_DEVMODE") == "1" {
		return http.Dir("./static")
	}
	return http.FS(assets)
}

func GetHandler() http.Handler {
	return http.FileServer(GetFilesystem())
}
