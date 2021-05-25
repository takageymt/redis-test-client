package main

import (
    "testing"
    "time"
    "net"
)

func TestPing(t *testing.T) {
    tests := []struct{
        name string
        input string
        expected string
    }{
        // PING
        {"PING", "PING", "+PONG\r\n"},
        {"PING pAQ94lNVFO", "PING pAQ94lNVFO", "$10\r\npAQ94lNVFO\r\n"},
        //{"PING pAQ94lNVFO", "PING", "pAQ94lNVFO"}, "+npAQ94lNVFO\r\n"},
    }
    tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:6379")
    conn, _ := net.DialTCP("tcp", nil, tcpAddr)
    conn.SetReadDeadline(time.Now().Add(5*time.Second))
    for _, testCase := range tests {
        actual, _ := Request(conn, testCase.input)
        if actual != testCase.expected {
            t.Errorf("%s: got unexpected reply: expected: %#v actual: %#v\n", testCase.name, testCase.expected, actual)
        }
    }
}

func TestGetSet(t *testing.T) {
    tests := []struct{
        name string
        input string
        expected string
    }{
        // SET and GET
        {"GET O1Gqx - 1", "GET O1Gqx", "$-1\r\n"}, // GETもSETもされていない時はnil
        {"GET O1Gqx - 2", "GET O1Gqx", "$-1\r\n"}, // GETをされたあともSET前ならnilのまま
        {"SET O1Gqx TwGGHNSfAm", "SET O1Gqx TwGGHNSfAm", "+OK\r\n"}, // SET
        {"GET O1Gqx - 3", "GET O1Gqx", "$10\r\nTwGGHNSfAm\r\n"}, // SETされたあとはGETできる
        {"SET O1Gqx mOsfXgDmXr", "SET O1Gqx mOsfXgDmXr", "+OK\r\n"}, // SETは上書き
        {"GET O1Gqx - 4", "GET O1Gqx", "$10\r\nmOsfXgDmXr\r\n"}, // SETは上書き
        {"SET aYEL2 MO/aO77/YF", "SET aYEL2 MO/aO77/YF NX", "+OK\r\n"}, // GETより前のSETは当然有効
    }
    tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:6379")
    conn, _ := net.DialTCP("tcp", nil, tcpAddr)
    conn.SetReadDeadline(time.Now().Add(5*time.Second))
    for _, testCase := range tests {
        actual, _ := Request(conn, testCase.input)
        if actual != testCase.expected {
            t.Errorf("%s: got unexpected reply: expected: %#v actual: %#v\n", testCase.name, testCase.expected, actual)
        }
    }
}

func TestSetNxXx(t *testing.T) {
    tests := []struct{
        name string
        input string
        expected string
    }{
        // SET NX
        {"SET O1Gqx CKP200y/Wz NX", "SET O1Gqx CKP200y/Wz NX", "$-1\r\n"},  // 存在する場合は上書きしない
        {"GET O1Gqx - 5", "GET O1Gqx", "$10\r\nmOsfXgDmXr\r\n"},
        {"SET 9+uva 7yny7kGspT NX", "SET 9+uva 7yny7kGspT NX", "+OK\r\n"}, // 存在しない場合はSETする
        {"GET 9+uva - 6", "GET 9+uva", "$10\r\n7yny7kGspT\r\n"},
        // SET XX
        {"SET O1Gqx SEGmfHG7tK XX", "SET O1Gqx SEGmfHG7tK XX", "+OK\r\n"}, // 存在する場合は上書きする
        {"GET O1Gqx - 3", "GET O1Gqx", "$10\r\nSEGmfHG7tK\r\n"}, 
        {"SET YygOT he+QoeLnkJ XX", "SET YygOT he+QoeLnkJ XX", "$-1\r\n"}, // 存在しない場合はSETしない
        {"GET YygOT - 4", "GET YygOT", "$-1\r\n"},
    }
    tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:6379")
    conn, _ := net.DialTCP("tcp", nil, tcpAddr)
    conn.SetReadDeadline(time.Now().Add(5*time.Second))
    for _, testCase := range tests {
        actual, _ := Request(conn, testCase.input)
        if actual != testCase.expected {
            t.Errorf("%s: got unexpected reply: expected: %#v actual: %#v\n", testCase.name, testCase.expected, actual)
        }
    }
}

func TestDel(t *testing.T) {
    tests := []struct{
        name string
        input string
        expected string
    }{
        // DEL
        {"DEL O1Gqx", "DEL O1Gqx", ":1\r\n"}, // 存在する場合は削除
        {"GET O1Gqx", "GET O1Gqx", "$-1\r\n"}, // 削除されたらGETはnil
        {"DEL oV7Q+", "DEL oV7Q+", ":0\r\n"}, // 存在しない場合は何もしない
        {"DEL YygOT", "DEL YygOT", ":0\r\n"}, // SET XX が失敗したら存在しないので何もしない
        {"DEL aYEL2", "DEL aYEL2", ":1\r\n"}, // GETされていなくても削除可能
        {"SET 05Zk+ dLxiXJEkOI", "SET 05Zk+ dLxiXJEkOI", "+OK\r\n"}, // SET
        {"SET 1ruFZ IGMydT37Oy", "SET 1ruFZ IGMydT37Oy", "+OK\r\n"}, // SETは上書き
        {"DEL 05Zk+ 1ruFZ", "DEL 05Zk+ 1ruFZ", ":2\r\n"}, // 複数削除も可能
        {"DEL 05Zk+ 9+uva 1ruFZ", "DEL 05Zk+ 9+uva 1ruFZ", ":1\r\n"}, // 存在しない:+0, 存在する:+1
    }
    tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:6379")
    conn, _ := net.DialTCP("tcp", nil, tcpAddr)
    conn.SetReadDeadline(time.Now().Add(5*time.Second))
    for _, testCase := range tests {
        actual, _ := Request(conn, testCase.input)
        if actual != testCase.expected {
            t.Errorf("%s: got unexpected reply: expected: %#v actual: %#v\n", testCase.name, testCase.expected, actual)
        }
    }
}

func TestIncrby(t *testing.T) {
    tests := []struct{
        name string
        input string
        expected string
    }{
        // INCRBY
        {"SET GlirZDSqMf 4613", "SET GlirZDSqMf 4613", "+OK\r\n"}, // 整数をSET
        {"INCRBY GlirZDSqMf 1934", "INCRBY GlirZDSqMf 1934", ":6547\r\n"}, // 正数を加算
        {"INCRBY GlirZDSqMf -6115", "INCRBY GlirZDSqMf -6115", ":432\r\n"}, // 負数を加算
        {"INCRBY GlirZDSqMf -4224", "INCRBY GlirZDSqMf -4224", ":-3792\r\n"}, // メモリの値を負にする
        {"INCRBY 7NDUJgUqSL 358", "INCRBY 7NDUJgUqSL 358", ":358\r\n"}, // 存在しないキーの値は0とする
        {"DEL GlirZDSqMf 7NDUJgUqSL", "DEL GlirZDSqMf 7NDUJgUqSL", ":2\r\n"},
    }
    tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:6379")
    conn, _ := net.DialTCP("tcp", nil, tcpAddr)
    conn.SetReadDeadline(time.Now().Add(5*time.Second))
    for _, testCase := range tests {
        actual, _ := Request(conn, testCase.input)
        if actual != testCase.expected {
            t.Errorf("%s: got unexpected reply: expected: %#v actual: %#v\n", testCase.name, testCase.expected, actual)
        }
    }
}
