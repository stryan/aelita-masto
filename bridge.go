package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func connectToAelita(host string, port string) net.Conn {
	c, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("sending header")
	fmt.Fprintf(c, "aelita 0.1\n")
	resp, _ := bufio.NewReader(c).ReadString('\n')
	if strings.TrimSpace(resp) == "OK aelita 0.1" {
		fmt.Println("connection established")
		return c
	}
	fmt.Println("Bad connection")
	return nil
}

func sendCommand(c net.Conn, cmd string) string {
	fmt.Fprintf(c, strings.TrimSpace(cmd)+"\n")
	resp, _ := bufio.NewReader(c).ReadString('\n')
	return resp
}

func disconnectFromAelita(c net.Conn) {
	fmt.Println("Closing connection")
	fmt.Fprintf(c, "close\n")
	fmt.Println("Check for clean reponse")
	resp, _ := bufio.NewReader(c).ReadString('\n')
	if strings.TrimSpace(resp) == "END" {
		fmt.Println("connection closed cleanly")
        } else {
		fmt.Println("connection closed poorly")
	}
	c.Close()
}
