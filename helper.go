package main

import (
	"github.com/mattn/go-runewidth"
    "github.com/acarl005/stripansi"
)

func calcWidhtColorRed(s string) int {
	return runewidth.StringWidth(stripansi.Strip(s))
}
