package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strings"
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

	u, err := url.Parse(*urlFlag)
	checkError(err)
	fmt.Println(u.Host)

	// "http://cloudflare-assignment.rahulnakre.workers.dev"
	// "http://www.google.com"
	// "http://ran-home.s3-website-us-east-1.amazonaws.com"

	var slowestResTime, fastestResTime float64 = -1, -1
	var largestResSize, smallestResSize int64 = -1, -1
	var totalTime float64 = 0
	var successCount int64 = 0
	// var totalBytesRead int =
	for i := int64(0); i < *profileFlag; i++ {
		startTime := time.Now()

		conn, err := net.Dial("tcp", u.Host+":http")
		checkError(err)
		defer conn.Close()

		// _, err = fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", u.Host)
		_, err = fmt.Fprintf(conn, "GET %s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", u.Path, u.Host)
		checkError(err)

		// status, err := bufio.NewReader(conn).ReadString('\n')
		// checkError(err)
		// fmt.Printf(status)

		_, totalBytesRead, err := Read(conn)
		// res, err := ioutil.ReadAll(conn)
		checkError(err)

		endTime := time.Now()

		// fmt.Printf(s + "\n")
		fmt.Printf("bytes read: %d\n", totalBytesRead)
		// fmt.Println(string(res) + "\n")

		// fmt.Printf("size %d\n", len(res))
		// if strings.Contains(status, "200") {
		// 	successCount++
		// }

		resTime := endTime.Sub(startTime).Seconds()
		if i == 0 || resTime < fastestResTime {
			fastestResTime = resTime
		}
		if i == 0 || resTime > slowestResTime {
			slowestResTime = resTime
		}
		totalTime += resTime

		if i == 0 || totalBytesRead < smallestResSize {
			smallestResSize = totalBytesRead
		}
		if i == 0 || totalBytesRead > largestResSize {
			largestResSize = totalBytesRead
		}
	}

	fmt.Printf("Number of requests: %d\n", *profileFlag)
	fmt.Printf("Fastest Response Time: %f\n", fastestResTime)
	fmt.Printf("Slowest Response Time: %f\n", slowestResTime)
	fmt.Printf("Mean Response Time: %f\n", float64(totalTime)/float64(*profileFlag))
	fmt.Printf("Median Response Time: %f\n", 1.0)
	fmt.Printf("Percentage of Successful Requests: %f%%\n", float64(successCount)/float64(*profileFlag)*100)
	fmt.Printf("Size of smallest response (in bytes): %d\n", smallestResSize)
	fmt.Printf("Size of largest response (in bytes): %d\n", largestResSize)
	// fmt.Printf("")
	// fmt.Println(u)

}

// Read reads from a connection
func Read(conn net.Conn) (string, int64, error) {
	reader := bufio.NewReader(conn)
	// scanner := bufio.NewScanner(reader)
	var buffer bytes.Buffer
	var i int = -1
	var totalBytesRead int64 = 0
	var headerDone bool = false
	// for scanner.Scan() {
	// 	// fmt.Println("here")
	// 	w := scanner.Text()
	// 	buffer.WriteString(w)
	// 	fmt.Println(w)
	// 	fmt.Printf("%v\n", strings.Contains(w, "\n"))
	// 	if len(w) >= 4 && w[len(w)-4:] == "\r\n\r\n" {
	// 		buffer.WriteString("\n")
	// 		fmt.Printf(buffer.String())
	// 		buffer.Reset()
	// 	}
	// }

	// return "sup", 3, nil
	var isPrevDelim bool = false
	for {
		i++
		// bytesArr, _, err := reader.ReadLine()
		bytesArr, err := reader.ReadBytes('\n')
		// if err != nil {

		// }
		if headerDone {
			totalBytesRead += int64(len(bytesArr))
		}
		// fmt.Println(i)
		// fmt.Println("ere")

		fmt.Printf(string(bytesArr))
		// fmt.Printf("bytes: %d\n", len(bytesArr))
		if err != nil {
			fmt.Println("didnt end in delim")
		}
		// fmt.Println(strings.Contains(string(bytesArr), "\r\n"))
		// fmt.Println(strings.EqualFold(string(bytesArr), "\r\n"))
		if !headerDone {
			if strings.EqualFold(string(bytesArr), "\r\n") && isPrevDelim {
				fmt.Println("END OF HEADER")
				headerDone = true
				continue
			}

			if strings.Contains(string(bytesArr), "\r\n") {
				isPrevDelim = true
			}
		} else {
			fmt.Printf("bytes: %d\n", len(bytesArr))
		}

		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
				break
			}
			return "", totalBytesRead, err
		}
		buffer.Write(bytesArr)
		// if headerDone {
		// 	totalBytesRead += len(bytesArr)
		// }

		// if !isPrefix {
		// 	fmt.Println("breakubg cus isPrefix")
		// 	break
		// }
	}
	// return buffer.String(), buffer.Len(), nil
	return buffer.String(), totalBytesRead, nil
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(0)
	}
}
