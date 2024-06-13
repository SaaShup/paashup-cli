package main

import (
	"github.com/mattn/go-runewidth"
	"strings"
)

func calcWidhtColorRed(s string) int {
	return runewidth.StringWidth(strings.Replace(strings.Replace(s, "\x1b[31m", "", 1), "\x1b[0m", "", 1))
}
