package main

import (
    "log"
    "net"
    "fmt"
    "strings"
    "strconv"
)

func encode(data []string) string {
    ret := make([]string, len(data)*2+2)
    ret[0] = "*"+strconv.Itoa(len(data))
    for i, s := range data {
        ret[i*2+1] = "$"+strconv.Itoa(len(s))
        ret[i*2+2] = s
    }
    ret[len(data)*2+1] = ""
    return strings.Join(ret, "\r\n")
}

func Request(conn *net.TCPConn, args []string) (string, error) {
    conn.Write([]byte(encode(args)))
    buf := make([]byte, 1024)
    n, _ := conn.Read(buf)
    return string(buf[:n]), nil
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Fatal("ResolveAddr", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal("DialTCP", err)
	}

    reply, err := Request(conn, []string{"PING"})
    fmt.Println(reply)
    reply, err = Request(conn, []string{"PING Hello"})
    fmt.Println(reply)
    reply, err = Request(conn, []string{"GET mykey"})
    fmt.Println(reply)
    reply, err = Request(conn, []string{"SET mykey value"})
    fmt.Println(reply)
    reply, err = Request(conn, []string{"GET mykey"})
    fmt.Println(reply)
    reply, err = Request(conn, []string{"DEL mykey"})
    fmt.Println(reply)
}
