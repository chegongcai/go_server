package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

//182.254.185.142  8080
func main() {
	//server
	service := ":8080"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErr(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)

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
		fmt.Println("rev data for go client: ", string(buf[0:n]))
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
	var err error

	arr_buf = strings.Split(rev_buf, "#") //先分割#

	serial_num := string(arr_buf[3])
	imei := string(arr_buf[2])
	ip := string(arr_buf[0])
	time.Sleep(5 * time.Second)
	//send data
	buf := fmt.Sprintf("%s#S168#%s#%s#0009#ACK^LOCA,$", ip, imei, serial_num)
	fmt.Println("send data to go client: ", buf)
	_, err = conn.Write([]byte(buf))
	if err != nil {
		return
	}
	fmt.Println("****************************************************************************************")
}

/*
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	//建立socket端口监听
	netListen, err := net.Listen("tcp", "localhost:1024")
	CheckError(err)

	defer netListen.Close()

	Log("Waiting for clients ...")

	//等待客户端访问
	for {
		conn, err := netListen.Accept() //监听接收
		if err != nil {
			continue //如果发生错误，继续下一个循环。
		}

		Log(conn.RemoteAddr().String(), "tcp connect success") //tcp连接成功
		go handleConnection(conn)
	}
}

//处理连接
func handleConnection(conn net.Conn) {
	buffer := make([]byte, 2048) //建立一个slice
	for {
		n, err := conn.Read(buffer) //读取客户端传来的内容
		if err != nil {
			Log(conn.RemoteAddr().String(), "connection error: ", err)
			return //当远程客户端连接发生错误（断开）后，终止此协程。
		}

		Log(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))

		//返回给客户端的信息
		//strTemp := "CofoxServer got msg \"" + string(buffer[:n]) + "\" at " + time.Now().String()
		strTemp := time.Now().String() + "Server have get the data\n"
		conn.Write([]byte(strTemp))
	}
}

//日志处理
func Log(v ...interface{}) {
	log.Println(v...)
}

//错误处理
func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
	}
}
*/
