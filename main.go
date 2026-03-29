package main

import (
	"embed"
	"os"

	"github.com/mockzilla/connexions/v2/pkg/portable"
)

//go:embed openapi
var content embed.FS

func main() {
	os.Exit(portable.RunFS(content, os.Args[1:]))
}
