package clients

import (
	"errors"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Clients struct {
	// Used for concurrent access to map and to ensure only one WriteJSON is active.
	// (Probably should use channels, but this works for now.)
	sync.RWMutex
	connections map[string]*Client
}

func New() *Clients {
	clients := &Clients{
		connections: make(map[string]*Client),
	}
	http.HandleFunc("/ws", clients.wsHandler)
	return clients
}

func (c *Clients) Names() []string {
	c.RLock()
	defer c.RUnlock()
	var a []string
	for name := range c.connections {
		a = append(a, name)
	}
	return a
}

func (c *Clients) Send(name string, message interface{}) error {
	c.RLock()
	defer c.RUnlock()

	client, ok := c.connections[name]
	if !ok {
		return errors.New("client not found or connected")
	}

	return client.activeConn.WriteJSON(message)
}

type Client struct {
	activeConn *websocket.Conn
}

type clientCommand struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

var upgrader = websocket.Upgrader{}

func (c *Clients) wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}

	defer ws.Close()

	client := Client{
		activeConn: ws,
	}
	name := ""

	for {
		var cmd clientCommand
		err := ws.ReadJSON(&cmd)
		if err != nil {
			if err == io.EOF {
			} else {
				log.Printf("Error reading from client: %v", err)
			}
			if len(name) > 0 {
				c.Lock()
				if c.connections[name] == &client {
					delete(c.connections, name)
				}
				c.Unlock()
			}
			break
		}

		switch cmd.Type {
		case "connect":
			if len(name) == 0 {
				name = cmd.Name
				c.Lock()
				c.connections[name] = &client
				c.Unlock()
			} else {
				log.Printf("Ignored duplicate connect from client: %#v", cmd)
			}
		case "ping":
		default:
			log.Printf("Unknown command from client: %#v", cmd)
		}
	}
}
