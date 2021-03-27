package static

import "embed"

//go:embed img/* js/* css/* fonts/*
//go:embed index.html
var Static embed.FS