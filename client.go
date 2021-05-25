package main

import (
    "log"
    "net"
    "fmt"
    "strings"
    "strconv"
)

func encode(cmd string) string {
    data := strings.Fields(cmd)
    ret := make([]string, len(data)*2+2)
    ret[0] = "*"+strconv.Itoa(len(data))
    for i, s := range data {
        ret[i*2+1] = "$"+strconv.Itoa(len(s))
        ret[i*2+2] = s
    }
    ret[len(data)*2+1] = ""
    return strings.Join(ret, "\r\n")
}

func Request(conn *net.TCPConn, args string) (string, error) {
    conn.Write([]byte(encode(args)))
    buf := make([]byte, 1024)
    n, _ := conn.Read(buf)
    return string(buf[:n]), nil
}
/*
func Ping(conn *net.TCPConn, arg interface{}) (string, error) {
    pingArgs := []string{"PING"}
    if s, ok := arg.(string); ok {
        pingArgs = append(pingArgs, s)
    }
    return Request(conn, pingArgs)
}

func Get(conn *net.TCPConn, key string) (string, error) {
    getArgs := []string{"GET", key}
    return Request(conn, getArgs)
}

func Set(conn *net.TCPConn, key string, value string, option interface{}) (string, error) {
    setArgs := []string{"SET", key, value}
    if s, ok := option.(string); ok {
        setArgs = append(setArgs, s)
    }
    return Request(conn, setArgs)
}

func Del(conn *net.TCPConn, keys []string) (string, error) {
    delArgs := []string{"DEL"}
    delArgs = append(delArgs, keys...)
    return Request(conn, delArgs)
}

func Incrby(conn *net.TCPConn, key string, value string) (string, error) {
    incrbyArgs := []string{"INCRBY", key, value}
    return Request(conn, incrbyArgs)
}
*/
func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Fatal("ResolveAddr", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal("DialTCP", err)
	}

    reply, err := Request(conn, "PING")
    fmt.Println(reply)
    reply, err = Request(conn, "PING Hello")
    fmt.Println(reply)
    reply, err = Request(conn, "GET mykey")
    fmt.Println(reply)
    reply, err = Request(conn, "SET mykey value")
    fmt.Println(reply)
    reply, err = Request(conn, "GET mykey")
    fmt.Println(reply)
    reply, err = Request(conn, "DEL mykey")
    fmt.Println(reply)
}
