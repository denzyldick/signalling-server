package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	var c Command;
  json.NewDecoder(r.Body).Decode(&c)
  if len(c.Command) > 0 {
  fmt.Println(c.Command,"Testing",len(c.Command))
	c.Device.Connection = w
	c.Execute()
	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}

	reader(ws)
}
}

func routes() {
	http.HandleFunc("/ws", wsEndpoint)
}

type Device struct {
	id         string
	Devices    []Device
	Connection http.ResponseWriter
}

func (d Device) addDevice(device Device) {
	d.Devices = append(d.Devices, device)
}

type Command struct {
  Command string `json:"command"`
  Device  Device `json:"device"`
	Sdp     string `json:"sdp"`
  Payload string `json:"payload"`
}

var Commands []Command

func (c Command) Execute() {
  fmt.Println("Executing command");

  fmt.Println(c.Command, "Command")
	if c.Command == "connect_devices" {
		Commands = append(Commands, c)
		fmt.Println("Starting new network")
	}

  if c.Command == "join_device" {
		fmt.Println("Joing a network")
		code := c.Payload

		for _, current := range Commands {
      fmt.Println(current.Device.id, code)
			if current.Device.id == code {
				current.Device.Connection.Write([]byte("Hello"))
				current.Device.Devices = append(current.Device.Devices, c.Device)
			}
		}
	}
}

func main() {
	routes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
