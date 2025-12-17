package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ahmedjebari022/go-docs/internal/api"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)



var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool{
		return true
	},
}

type Client struct{
	id    uuid.UUID 
	documentId  string
	conn *websocket.Conn
	sent	chan []Operation
	hub *Hub
}

type Message struct{
	UserId  	 uuid.UUID `json:"user_id"`
	DocumentId 	 string	 `json:"document_id"`
	Op 			 []Operation `json:"operation"`
}

type Hub struct{
	clients map[*Client]bool
	subscribe chan *Client
	unsubscribe chan *Client
	broadcast chan Message

}


func NewHub() Hub {
	return Hub{
		clients: make(map[*Client]bool),
		subscribe: make(chan *Client),
		unsubscribe: make(chan *Client),
		broadcast: make(chan Message),
	}
}


func (h *Hub) Run() {
	for {
		select {
		case client := <- h.subscribe:
			fmt.Printf("subscribing client :%v\n", client)
			h.clients[client] = true
			fmt.Printf("clients after subscribe: %v\n", h.clients)  // Add this
		case client := <- h.unsubscribe:
			if _, ok := h.clients[client]; ok {
				client.conn.Close()
				delete(h.clients,client)
			}
		case msg := <- h.broadcast:
			fmt.Printf("received brodcast :%v\n",msg)
			for c, _ := range h.clients {
				if c.documentId == msg.DocumentId && c.id != msg.UserId{
					c.sent <- msg.Op
				}
			}
			}
		}	
}
type Operation struct{
	Type  string  `json:"type"`
	Path  []int   `json:"path"`
	Offset int  `json:"offset"`
	Text	string `json:"text"`
}
type Operations struct{
	Ops		[]Operation  `json:"operations"`
}

func (c *Client) Reader(){
	defer func(){
		c.hub.unsubscribe <- c
		c.conn.Close()
	}()
	for {
		fmt.Println("Reader waiting for message...")
		_, reader, err := c.conn.NextReader()
		fmt.Println("Reader got message!")
		if err != nil {
			c.hub.unsubscribe <- c
			break
		}
		decoder := json.NewDecoder(reader)
		var op []Operation
		if err := decoder.Decode(&op); err != nil {
			fmt.Printf("error while deconding the json msg :%s\n", err.Error())
			break
		}
		fmt.Printf("Read :%v\n",op)
		
		c.hub.broadcast <- Message{
			DocumentId: c.documentId,
			Op: op,
			UserId: c.id,
		}
	}
}

func (c *Client) Writer(){
	defer func(){
		c.hub.unsubscribe <- c
		c.conn.Close()
	}()
	for op := range c.sent{
		fmt.Println("client got broadcast")
		writer, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			c.hub.unsubscribe <- c
		}
		fmt.Printf("doc: %v",op)
		encoder := json.NewEncoder(writer)
		if err := encoder.Encode(op); err != nil {
			fmt.Printf("error while encodin the doc :%s\n",err.Error())
		}
		if err := writer.Close(); err != nil {
			return
		}
	}
}

func (h *Hub)wsHandler(w http.ResponseWriter, r *http.Request){
	userId, err := api.GetUserIdFromContext(r.Context())
	if err != nil {
		fmt.Println(err.Error())
		return 
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	fmt.Println("Connecting to websocket")
	if err != nil {
		fmt.Println("Problem upgrading")
		return
	}		
	documentIdString := r.PathValue("documentId")	
	_, err = uuid.Parse(documentIdString)
	if err != nil {
		fmt.Printf("error parsin the document id :%s\n",err.Error())
	}
	c := &Client{
		id: userId,
		documentId: documentIdString,
		conn: conn,
		hub: h,
		sent: make(chan []Operation),
	}
	h.subscribe <- c
	go c.Reader()
	go c.Writer()
}
