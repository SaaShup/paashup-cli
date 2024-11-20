package utils

import (
	"github.com/mattn/go-runewidth"
    "github.com/acarl005/stripansi"
)

func CalcWidhtColorRed(s string) int {
	return runewidth.StringWidth(stripansi.Strip(s))
}
