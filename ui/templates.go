package ui

import (
	"html/template"
	"io/ioutil"

	"github.com/dgl/go-web-dashboard/web"
)

var rootTemplate = template.Must(template.New("root").Parse(mustReadAll("index.html")))
var showTemplate = template.Must(template.New("show").Parse(mustReadAll("show.html")))

func mustReadAll(file string) string {
	f, err := web.Template.Open(file)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(b)
}
