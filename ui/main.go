package ui

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dgl/go-web-dashboard/clients"
	"github.com/dgl/go-web-dashboard/web"
)

type UI struct {
	clients *clients.Clients
}

func New(clients *clients.Clients) *UI {
	ui := UI{
		clients: clients,
	}

	http.HandleFunc("/", ui.root)
	http.HandleFunc("/show", ui.show)
	http.HandleFunc("/blank", serveFile("blank.html"))
	http.HandleFunc("/time", serveFile("time.html"))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(web.JS)))
	http.HandleFunc("/send", ui.send)

	return &ui
}

func serveFile(file string) func(w http.ResponseWriter, r *http.Request) {
	content := mustReadAll(file)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(content))
	}
}

func (ui *UI) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	err := rootTemplate.Execute(w, map[string]interface{}{
		"Clients": ui.clients.Names(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ui *UI) show(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if len(name) == 0 {
		cookie, err := r.Cookie("gww_name")
		if err == nil {
			name = cookie.Value
		}
	}
	err := showTemplate.Execute(w, map[string]interface{}{
		"Name": name,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ui *UI) send(w http.ResponseWriter, r *http.Request) {
	// Either form data or a direct post of JSON...
	name := r.FormValue("name")
	var data interface{}

	if r.Header.Get("Content-Type") == "application/json" {
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(name) == 0 {
			if d, ok := data.(map[string]interface{}); ok {
				name = d["name"].(string)
			}
		}
	} else {
		url := r.FormValue("url")
		data = map[string]string{"url": url}
	}

	if len(name) > 0 {
		err := ui.clients.Send(name, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write([]byte("ok"))
	} else {
		log.Printf("Missing name, ignored /send")
	}
}
