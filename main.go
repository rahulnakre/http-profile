package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/pborman/getopt"
)

func main() {
	// https://cloudflare-assignment.rahulnakre.workers.dev/
	urlFlag := getopt.StringLong("url", 0, "", "url of a site")
	helpFlag := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

	if *helpFlag {
		getopt.Usage()
		os.Exit(0)
	}

	fmt.Println("sup " + *urlFlag + " hehe")

	tcpAddr, err := net.ResolveTCPAddr("tcp4", "www.google.com")
	errHandler(err)

	fmt.Printf("here: %v", tcpAddr)

	conn, err := net.Dial("tcp", *urlFlag)
	errHandler(err)
	defer conn.Close()

	_, err = conn.Write([]byte("GET / HTTP/1.1\r\n\r\n"))
	errHandler(err)

	res, err := ioutil.ReadAll(conn)
	errHandler(err)

	fmt.Println(string(res))

	// buf := make([]byte, 0, 4096)
	// recv := make([]byte, 256)
	// for {
	// 	fmt.Println("fere")
	// 	size, err := conn.Read(recv)
	// 	fmt.Printf("read size: %d", size)

	// 	// errHandler(err)
	// 	if err != nil {
	// 		break
	// 	}
	// 	buf = append(buf, recv[:size]...)
	// }
	// fmt.Printf("%s", string(buf))
	// conn.

	// addr, err := net.LookupAddr("ispycode.com")
	// errHandler(err)

	// fmt.Println(addr)

	// tcpAddr, err := net.ResolveTCPAddr("tcp", "8.8.8.8")
	// errHandler(err)
	// net.DialIP()
	// fmt.Println(tcpAddr)
	// conn, err := net.DialTCP("tcp", nil, tcpAddr)
	// errHandler(err)

	// fmt.Printf("%v\n", conn)

}

func errHandler(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(0)
	}
}
