package main

import (
	"bufio"
	"fmt"
	"net"

	errors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	const (
		address = "127.0.0.1:8080"
	)

	log.SetLevel(log.DebugLevel)
	log.Debugf("Starting tcp on %v", address)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Error(err.Error())
	}

	aconns := make(map[net.Conn]int)
	conns := make(chan net.Conn)
	dconns := make(chan net.Conn)
	msgs := make(chan string)
	i := 0

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Error(err.Error())
			}
			conns <- conn
		}
	}()

	for {
		select {
		case conn := <-conns:
			aconns[conn] = i
			i++
			go func(conn net.Conn, i int) {
				rd := bufio.NewReader(conn)
				for {
					m, err := rd.ReadString('\n')
					if err != nil {
						log.Debug(errors.Wrap(err, "Done reading").Error())
						break
					}
					msgs <- fmt.Sprintf("Client %v: %v", i, m)
				}
				dconns <- conn
			}(conn, i)
		case msg := <-msgs:
			for conn := range aconns {
				conn.Write([]byte(msg))
			}

		case dconn := <-dconns:
			log.Infof("Client %v has disconnected", aconns[dconn])
			delete(aconns, dconn)

		}
	}

}
