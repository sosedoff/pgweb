package helpers

import (
	"fmt"

	"github.com/sclevine/agouti"
)


func FillConnectionForm(page *agouti.Page, data map[string]string) {
	for selector, value := range data {
		page.Find(selector).Fill(value)
	}
}



func Screenshot(page *agouti.Page, name string) {
	page.Screenshot(fmt.Sprintf("_output/%s.png", name))
}

