package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"ssgb-matching/messages"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ServerUrl string
}

type ClientTicket struct {
	Id string
}

type ClientConnection struct {
	Conn *websocket.Conn
}

func parseFlags() (Client, int) {
	serverUrl := flag.String("s", "127.0.0.1:9990", "server url")
	n := flag.Int("n", 1, "num requests")

	flag.Parse()
	return Client{
		ServerUrl: *serverUrl,
	}, *n
}

func (c *Client) getTicket() (ClientTicket, error) {
	t := ClientTicket{}
	form := url.Values{
		"class": {"1"},
	}
	serverUrl := "http://" + c.ServerUrl + "/ticket/new"
	res, err := http.PostForm(serverUrl, form)
	if err != nil {
		return t, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return t, err
	}

	err = json.Unmarshal(b, &t)
	return t, err
}

func (c *Client) startListenTicket(ticket ClientTicket) (ClientConnection, error) {
	cc := ClientConnection{}
	serverUrl := "ws://" + c.ServerUrl + "/ticket/listen/" + ticket.Id
	conn, res, err := websocket.DefaultDialer.Dial(serverUrl, nil)
	if err != nil {
		return cc, err
	} else if res.StatusCode != http.StatusSwitchingProtocols {
		return cc, errors.New("protocol switching does not work")
	}
	defer res.Body.Close()

	cc.Conn = conn
	return cc, nil
}

func (c *ClientConnection) listenTicket() {
	for {
		msg := messages.StatusMessage{}
		if err := c.Conn.ReadJSON(&msg); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", msg)
		if msg.Status == messages.StatusMatched ||
			msg.Status == messages.StatusError {

			break
		}
	}
}

func (c *Client) processMatching(onDone func()) {
	ticket, err := c.getTicket()
	if err != nil {
		panic(err)
	}

	fmt.Printf("ticket id: %s\n", ticket.Id)

	conn, err := c.startListenTicket(ticket)
	if err != nil {
		panic(err)
	}

	conn.listenTicket()

	onDone()
}

func main() {
	c, n := parseFlags()
	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go c.processMatching(func() {
			wg.Done()
		})
	}

	wg.Wait()
}
