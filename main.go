package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

var (
	fromHost  = flag.String("from-host", "127.0.0.1", "Listen host")
	fromPort = flag.Int("from-port", 0, "Listen port")
	toHost   = flag.String("to-host", "", "Forward to host")
	toPort   = flag.Int("to-port", 0, "Forward to port")
	verbose  = flag.Bool("v", false, "Verbose output")
)

func main() {
	flag.Parse()

	if *fromPort == 0 || *toHost == "" || *toPort == 0 {
		fmt.Println("PortForward - TCP Port Forwarder")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  portforward -from-port 8080 -to-host 127.0.0.1 -to-port 80")
		fmt.Println("")
		flag.PrintDefaults()
		return
	}

	addr := fmt.Sprintf("%s:%d", *fromHost, *fromPort)
	toAddr := fmt.Sprintf("%s:%d", *toHost, *toPort)

	fmt.Printf("🔄 PortForward\n")
	fmt.Printf("   Listening: %s\n", addr)
	fmt.Printf("   Forwarding: %s\n", toAddr)
	fmt.Printf("   Press Ctrl+C to stop\n")

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("❌ Failed to listen: %v\n", err)
		return
	}
	defer ln.Close()

	var wg sync.WaitGroup
	connID := 0

	for {
		src, err := ln.Accept()
		if err != nil {
			continue
		}

		connID++
		id := connID

		if *verbose {
			fmt.Printf("📥 Connection #%d from %s\n", id, src.RemoteAddr())
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			forward(src, toAddr, id)
		}()
	}
}

func forward(src net.Conn, toAddr string, id int) {
	defer src.Close()

	dst, err := net.DialTimeout("tcp", toAddr, 10*time.Second)
	if err != nil {
		if *verbose {
			fmt.Printf("❌ #%d: Failed to connect: %v\n", id, err)
		}
		return
	}
	defer dst.Close()

	if *verbose {
		fmt.Printf("📤 #%d: Connected to %s\n", id, toAddr)
	}

	// Bidirectional copy
	done := make(chan bool, 2)

	copy := func(w io.Writer, r io.Reader) {
		n, _ := io.Copy(w, r)
		if *verbose {
			fmt.Printf("📋 #%d: Transferred %d bytes\n", id, n)
		}
		done <- true
	}

	go copy(dst, src)
	go copy(src, dst)

	<-done
}