package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

///////////////// HUB /////////////////

type Hub struct {
	CentralClient         *Client
	PlayerClients         map[*Client]bool
	AddCentralClient      chan *Client
	AddPlayerClient       chan *Client
	RemoveClient          chan *Client
	PlayerClientBroadcast chan []byte
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.AddCentralClient:
			fmt.Println("adding central client to hub")
			h.CentralClient = client
		case client := <-h.AddPlayerClient:
			fmt.Println("adding player client to hub")
			h.PlayerClients[client] = true
		case client := <-h.RemoveClient:
			if client == h.CentralClient {
				fmt.Println("removing central client from hub")
				h.CentralClient = nil
			} else {
				fmt.Println("removing player client from hub")
				delete(h.PlayerClients, client)
			}
			client.Cleanup()
		case message := <-h.PlayerClientBroadcast:
			for c := range h.PlayerClients {
				c.Send <- message
			}
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		CentralClient:         nil,
		PlayerClients:         make(map[*Client]bool),
		AddCentralClient:      make(chan *Client),
		AddPlayerClient:       make(chan *Client),
		RemoveClient:          make(chan *Client),
		PlayerClientBroadcast: make(chan []byte),
	}
}

///////////////// CLIENT /////////////////

type Client struct {
	Hub *Hub
	// Ws       *websocket.Conn
	ClientId string
	Send     chan []byte
}

func (cl *Client) RecieveMessages() {
	// do player removal from game state and websocket close on disconnect
	defer func() {
		cl.HandleClientDisconnect(nil)
		// cl.Ws.Close()
	}()
	for {
		// // read message
		// _, message, err := cl.Ws.ReadMessage()
		// if err != nil {
		// 	log.Println("read:", err)
		// 	break
		// }
		// // log message received
		// // fmt.Println("client message received:")
		// // ConsoleLogJsonByteArray(message)
		// // route message to handler
		// messageTypeToHandler := map[string]func([]byte){
		// 	"MESSAGE_TYPE_PLAYER_JOIN":       cl.RouteMessageToCentralClient,
		// 	"MESSAGE_TYPE_PLAYER_LEAVE":      cl.RouteMessageToCentralClient,
		// 	"MESSAGE_TYPE_PLAYER_INPUT":      cl.RouteMessageToCentralClient,
		// 	"MESSAGE_TYPE_GAME_STATE":        cl.BroadcastMessageToPlayerClients,
		// 	"MESSAGE_TYPE_CLIENT_DISCONNECT": cl.HandleClientDisconnect,
		// }
		// var mData map[string]interface{}
		// if err := json.Unmarshal(message, &mData); err != nil {
		// 	panic(err)
		// }
		// // process message with handler
		// messageTypeToHandler[mData["messageType"].(string)](message)
	}
}

func (cl *Client) SendMessages() {
	// for message := range cl.Send {
	// 	cl.Ws.WriteMessage(1, message)
	// 	// log that message was sent
	// 	// fmt.Println("server message sent:")
	// 	// ConsoleLogJsonByteArray(messageJson)
	// }
}

func (cl *Client) RouteMessageToCentralClient(m []byte) {
	cl.Hub.CentralClient.Send <- m
}

func (cl *Client) BroadcastMessageToPlayerClients(m []byte) {
	cl.Hub.PlayerClientBroadcast <- m
}

func (cl *Client) HandleClientDisconnect(m []byte) {
	cl.Hub.RemoveClient <- cl
}

func (cl *Client) Cleanup() {
	close(cl.Send)
}

///////////////// RUN SERVER /////////////////

func main() {
	flag.Parse()
	log.SetFlags(0)
	// create and run hub singleton
	h := NewHub()
	go h.Run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello\n")
	})
	addr := flag.String("addr", "0.0.0.0:5000", "http service address")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

///////////////// HELPERS /////////////////

func ConsoleLogJsonByteArray(message []byte) {
	var out bytes.Buffer
	message = append(message, "\n"...)
	json.Indent(&out, message, "", "  ")
	out.WriteTo(os.Stdout)
}
