package web

// +build dev

import (
	"net/http"
)

var JS http.FileSystem = http.Dir("web/js")

var Template http.FileSystem = http.Dir("web/templates")
