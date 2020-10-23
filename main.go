package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pborman/getopt"
	myheap "github.com/rahulnakre/systems_assignment/heap"
)

func main() {
	// For the median
	minHeap := &myheap.MinHeap{}
	heap.Init(minHeap)
	maxHeap := &myheap.MaxHeap{}
	heap.Init(maxHeap)

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
	if len(u.Path) == 0 {
		u.Path = "/"
	}

	// "http://cloudflare-assignment.rahulnakre.workers.dev"
	// "http://www.google.com"
	// "http://ran-home.s3-website-us-east-1.amazonaws.com"

	var slowestResTime, fastestResTime float64 = -1, -1
	var largestResSize, smallestResSize int64 = -1, -1
	var totalTime float64 = 0
	var successCount int64 = 0
	var errorCodesArr []int
	var resTimeArr []float64
	// var totalBytesRead int =
	for i := int64(0); i < *profileFlag; i++ {
		startTime := time.Now()

		conn, err := net.Dial("tcp", u.Host+":http")
		checkError(err)
		defer conn.Close()

		_, err = fmt.Fprintf(conn, "GET %s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", u.Path, u.Host)
		checkError(err)

		_, totalBytesRead, statusCode, err := Read(conn)
		checkError(err)

		endTime := time.Now()

		fmt.Printf("bytes read: %d\n", totalBytesRead)

		if statusCode >= 200 && statusCode <= 299 {
			successCount++
		}

		if statusCode >= 400 && statusCode <= 599 {
			errorCodesArr = append(errorCodesArr, statusCode)
		}

		resTime := endTime.Sub(startTime).Seconds()

		resTimeArr = append(resTimeArr, resTime)

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

	sort.Float64s(resTimeArr)
	var median float64
	if len(resTimeArr)%2 == 0 {
		median = float64(resTimeArr[len(resTimeArr)/2]+resTimeArr[(len(resTimeArr)/2)-1]) / float64(2)
	} else {
		median = resTimeArr[len(resTimeArr)/2]
	}

	fmt.Printf("Number of requests: %d\n", *profileFlag)
	fmt.Printf("Fastest Response Time: %f\n", fastestResTime)
	fmt.Printf("Slowest Response Time: %f\n", slowestResTime)
	fmt.Printf("Mean Response Time: %f\n", float64(totalTime)/float64(*profileFlag))
	fmt.Printf("Median Response Time: %f\n", median)
	fmt.Printf("Percentage of Successful Requests: %f%%\n", float64(successCount)/float64(*profileFlag)*100)
	fmt.Printf("Error response codes: %v\n", errorCodesArr)
	fmt.Printf("Size of smallest response (in bytes): %d\n", smallestResSize)
	fmt.Printf("Size of largest response (in bytes): %d\n", largestResSize)
}

// Read reads from a connection
func Read(conn net.Conn) (string, int64, int, error) {
	reader := bufio.NewReader(conn)
	var buffer bytes.Buffer
	var i int = -1
	var statusCode int
	var totalBytesRead int64 = 0
	var headerDone bool = false
	var isPrevDelim bool = false

	for {
		i++
		bytesArr, err := reader.ReadBytes('\n')
		if headerDone {
			totalBytesRead += int64(len(bytesArr))
		}
		// totalBytesRead += int64(len(bytesArr))

		if err != nil {
			fmt.Println("didnt end in the delimiter")
			// return "", -1, -1, errors.New("Error reading a the request")
		}

		if !headerDone {
			if strings.EqualFold(string(bytesArr), "\r\n") && isPrevDelim {
				fmt.Printf("END OF HEADER: %d\n", totalBytesRead)
				totalBytesRead = 0
				headerDone = true
				continue
			}

			if strings.Contains(string(bytesArr), "\r\n") {
				isPrevDelim = true
			}

			if i == 0 {
				fmt.Println(strings.Split(string(bytesArr), " ")[1])
				statusCode, err = strconv.Atoi(strings.Split(string(bytesArr), " ")[1])
			}
		} else {
			fmt.Println(string(bytesArr))
			fmt.Printf("bytes: %d\n", len(bytesArr))
		}

		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
				break
			}
			return "", totalBytesRead, -1, err
		}
		buffer.Write(bytesArr)
	}
	return buffer.String(), totalBytesRead, statusCode, nil
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(0)
	}
}
