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
}

type Client struct{
	documentId  string
	conn *websocket.Conn
	sent	chan api.Document
	hub *Hub
}

type Message struct{
	Document 	api.Document `json:"document"`
	DocumentId 	 string	 `json:"document_id"`
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
				if c.documentId == msg.DocumentId{
					c.sent <- msg.Document
				}
			}
			}
		}	
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
		var doc api.Document
		if err := decoder.Decode(&doc); err != nil {
			fmt.Printf("error while deconding the json msg :%s\n", err.Error())
			break
		}
		fmt.Printf("Read :%v\n",doc)
		
		c.hub.broadcast <- Message{
			DocumentId: c.documentId,
			Document: doc,
		}
	}
}

func (c *Client) Writer(){
	defer func(){
		c.hub.unsubscribe <- c
		c.conn.Close()
	}()
	for doc := range c.sent{
		fmt.Println("client got broadcast")
		writer, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			c.hub.unsubscribe <- c
		}
		fmt.Printf("doc: %v",doc)
		encoder := json.NewEncoder(writer)
		if err := encoder.Encode(doc); err != nil {
			fmt.Printf("error while encodin the doc :%s\n",err.Error())
		}
		if err := writer.Close(); err != nil {
			return
		}
	}
}

func (h *Hub)wsHandler(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}		
	documentIdString := r.PathValue("documentId")	
	_, err = uuid.Parse(documentIdString)
	if err != nil {
		fmt.Println("error parsin the document id :%s",err.Error())
	}
	c := &Client{
		documentId: documentIdString,
		conn: conn,
		hub: h,
		sent: make(chan api.Document),
	}
	h.subscribe <- c
	go c.Reader()
	go c.Writer()
}
