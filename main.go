package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var client_conn net.Conn
var client_err error

func main() {
	//server
	service := ":8080" //182.254.185.142  8080
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErr(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)

	//client
	addr := "134.175.174.222:8080"
	client_conn, client_err = net.Dial("tcp", addr)
	if client_err != nil {
		log.Fatal(client_err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var buf [1024]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		rAddr := conn.RemoteAddr()
		fmt.Println("****************************************************************************************")
		fmt.Println("client ip: ", rAddr.String())
		fmt.Println("time: ", GetTimeStamp())
		fmt.Println("rev data: ", string(buf[0:n]))
		if buf[n-1] != '$' {
			return
		}
		rev_buf := string(buf[0 : n-1]) //delete the tail #
		ParseProtocol(rev_buf, conn)    //do protocol parse
	}
}

func GetTimeStamp() string {
	buf := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
	return buf
}

func ParseProtocol(rev_buf string, conn net.Conn) {
	var arr_buf []string

	arr_buf = strings.Split(rev_buf, "#") //先分割#

	serial_num := string(arr_buf[2])
	imei := string(arr_buf[1])

	//send data
	buf := fmt.Sprintf("S168#%s#%s#0009#ACK^LOCA,$", imei, serial_num)
	fmt.Println("send data: ", buf)
	_, client_err = client_conn.Write([]byte(buf))
	if client_err != nil {
		return
	}
	fmt.Println("****************************************************************************************")
}
