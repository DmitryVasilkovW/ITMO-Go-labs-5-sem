//go:build !solution

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var addr = flag.String("addr", "localhost:8080", "address")

func main() {
	flag.Parse()

	stop := setupSignalHandler()
	mes := readInput(os.Stdin)

	c, done := establishConnection(*addr)
	defer c.Close()

	handleMessages(c, done, stop, mes)
}

func setupSignalHandler() chan os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	return stop
}

func readInput(input *os.File) chan string {
	ch := make(chan string)
	go func() {
		scanner := bufio.NewScanner(input)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
		close(ch)
	}()
	return ch
}

func establishConnection(addr string) (*websocket.Conn, chan struct{}) {
	c, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}
			fmt.Print(string(message))
		}
	}()

	return c, done
}

func handleMessages(c *websocket.Conn, done chan struct{}, stop chan os.Signal, mes chan string) {
	for {
		select {
		case <-done:
			return
		case <-stop:
			log.Println("interrupt")
			gracefulClose(c, done)
			return
		case text := <-mes:
			if err := c.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}

func gracefulClose(c *websocket.Conn, done chan struct{}) {
	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("close write:", err)
		return
	}
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
	}
}
