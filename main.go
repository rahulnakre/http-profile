package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/pborman/getopt"
)

func main() {
	// https://cloudflare-assignment.rahulnakre.workers.dev/
	urlFlag := getopt.StringLong("url", 0, "", "url of a site")
	helpFlag := getopt.BoolLong("help", 0, "Help")
	profileFlag := getopt.Int64Long("profile", 0, 1, "Profile")
	getopt.Parse()

	if *helpFlag {
		getopt.Usage()
		os.Exit(0)
	}

	// tcpAddr, err := net.ResolveTCPAddr("tcp4", "www.google.com")
	// errHandler(err)

	// fmt.Printf("here: %v", tcpAddr)

	// conn, err := net.Dial("tcp", *urlFlag)
	// errHandler(err)
	// defer conn.Close()

	// _, err = conn.Write([]byte("GET / HTTP/1.1\r\n\r\n"))
	// errHandler(err)

	// res, err := ioutil.ReadAll(conn)
	// errHandler(err)

	// fmt.Println(string(res))

	// s := "http://ran-home.s3-website-us-east-1.amazonaws.com"

	u, err := url.Parse(*urlFlag)
	checkError(err)
	fmt.Println(u.Host)

	startTime := time.Now()
	// conn, err := net.Dial("tcp", "cloudflare-assignment.rahulnakre.workers.dev:http")
	// conn, err := net.Dial("tcp", u.Host+":http")
	conn, err := net.Dial("tcp", "www.google.com:http")
	// conn, err := net.Dial("tcp", "ran-home.s3-website-us-east-1.amazonaws.com:http")

	checkError(err)
	defer conn.Close()
	size, err := fmt.Fprintf(conn, "GET / HTTP/1.0\r\nHost: %s\r\nConnection: close\r\n\r\n", u.Host)
	checkError(err)
	fmt.Println(size)

	res, err := ioutil.ReadAll(conn)
	endTime := time.Now()
	checkError(err)

	fmt.Printf(string(res) + "\n")
	fmt.Printf("%f\n", endTime.Sub(startTime).Seconds())

	fmt.Println(*profileFlag)

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(0)
	}
}
