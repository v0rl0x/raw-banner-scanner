package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	bannerString string
	port         string
	output       string
	timeout      time.Duration
	threads      int
)

func init() {
	flag.StringVar(&bannerString, "b", "", "string to search for within the banner")
	flag.StringVar(&port, "p", "80", "port number")
	flag.StringVar(&output, "o", "output.txt", "output file")
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "timeout for banner grab")
	flag.IntVar(&threads, "t", 1, "number of threads")
}

func grabBanner(ip string, port string) (string, error) {
	conn, err := net.DialTimeout("tcp", ip+":"+port, timeout)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(timeout))

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer[:n]), nil
}

func main() {
	flag.Parse()

	var wg sync.WaitGroup
	sem := make(chan struct{}, threads)

	file, err := os.Create(output)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ip := scanner.Text()
		wg.Add(1)
		sem <- struct{}{}

		go func(ip string) {
			defer func() {
				<-sem
				wg.Done()
			}()

			banner, err := grabBanner(ip, port)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error grabbing banner: %v\n", err)
				return
			}

			if contains(banner, bannerString) || bannerString == "" {
				outputString := fmt.Sprintf("%s:%s\n%s\n\n", ip, port, banner)
				fmt.Fprint(writer, outputString)
				writer.Flush()
			}
		}(ip)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from input:", err)
	}

	wg.Wait()
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
