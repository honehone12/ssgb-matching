package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"ssgb-matching/gsip"
	"ssgb-matching/messages"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Bad       bool
	Backfill  bool
	ServerUrl string
}

type ClientTicket struct {
	Id string

	FoundBackfill bool
	BackfillGsip  gsip.GSIP
}

type ClientConnection struct {
	Conn *websocket.Conn
}

func parseFlags() (Client, int) {
	serverUrl := flag.String("s", "127.0.0.1:9990", "server url")
	n := flag.Int("n", 1, "num requests")
	backfill := flag.Bool("b", false, "use backfill")
	bad := flag.Bool("bad", false, "run bad client")

	flag.Parse()

	if *bad {
		return Client{
			ServerUrl: *serverUrl,
			Backfill:  *backfill,
			Bad:       true,
		}, 1
	}

	return Client{
		ServerUrl: *serverUrl,
		Backfill:  *backfill,
		Bad:       false,
	}, *n
}

func (c *Client) getTicket(class int64) (ClientTicket, error) {
	t := ClientTicket{}

	form := url.Values{
		"class":    {strconv.FormatInt(class, 10)},
		"backfill": {strconv.FormatBool(c.Backfill)},
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

	if res.StatusCode != http.StatusOK {
		return t, errors.New(string(b))
	}

	err = json.Unmarshal(b, &t)
	fmt.Printf("%#v\n", t)
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

func (c *ClientConnection) listenTicket(bad bool) {
	defer c.Conn.Close()

	for {
		if bad {
			n := rand.Intn(10)
			if n == 7 {
				fmt.Println("error 7, quitting")
				break
			}
		}

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

func (c *Client) processMatching(class int64, onDone func()) {
	ticket, err := c.getTicket(class)
	if err != nil {
		panic(err)
	}

	if ticket.FoundBackfill {
		onDone()
		return
	}

	conn, err := c.startListenTicket(ticket)
	if err != nil {
		panic(err)
	}

	conn.listenTicket(c.Bad)

	onDone()
}

const (
	classes = 6
)

func main() {
	c, n := parseFlags()
	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		class := rand.Int63n(classes) + 1
		go c.processMatching(class, func() {
			wg.Done()
		})
	}

	wg.Wait()
}
