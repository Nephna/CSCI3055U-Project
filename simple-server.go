package main

import (
    "fmt"
    "net"
    "os"
)

func main() {
    var message = make([]byte, 1024)

    service := ":52780"
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    checkError(err)
    listener, err := net.ListenTCP("tcp", tcpAddr)
    checkError(err)
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        _, err = conn.Read(message)
        checkError(err)
        conn.Write([]byte(message)
        conn.Close()
    }
}
func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}