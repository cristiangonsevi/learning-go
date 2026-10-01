package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

func startTunnel(client *ssh.Client, service Service, ready chan<- error) {
	localAddress := fmt.Sprintf("localhost:%d", service.LocalPort)
	listener, err := net.Listen("tcp", localAddress)
	if err != nil {
		log.Println("Error listener ", err)
		ready <- err
		return
	}
	defer listener.Close()
	ready <- nil
	for {

		localConn, err := listener.Accept()
		if err != nil {
			log.Println("Error local connections ", err)
			return
		}
		remoteAddres := fmt.Sprintf("localhost:%d", service.RemotePort)
		remoteConn, err := client.Dial("tcp", remoteAddres)
		if err != nil {
			log.Println("Error remote conncetion ", err)
			localConn.Close()
			continue
		}

		go func() {
			io.Copy(remoteConn, localConn)
			remoteConn.Close()
		}()
		go func() {
			io.Copy(localConn, remoteConn)
			localConn.Close()
		}()
	}
}
