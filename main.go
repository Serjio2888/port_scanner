package main

import (
	"runtime"
	"fmt"
	"sort"
	"net"
	"flag"
	"strconv"
)

var (
	WORKERS int
	PORTS string
	ADDRESS string
)

func worker(ports, results chan int) {
	for p := range ports {
		conn, err := net.Dial("tcp", ADDRESS + ":" + strconv.Itoa(p))
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}

}

func main() {
	flag.IntVar(&WORKERS, "c", WORKERS,"Number of CPU cores")
	flag.StringVar(&PORTS, "p", PORTS, "CHOSE YOUR PORTS: xxxx-xxxx")
	flag.StringVar(&ADDRESS, "a", ADDRESS, "URl of your site")
	flag.Parse()

	runtime.GOMAXPROCS(WORKERS)

	var char, char2 string
	for _, val := range PORTS {
		char = char + val
		if val == "-" {
			char2 = char
			char = ""
			continue
		}
	}

	var openports []int

	beg, err := strconv.Atoi(char2)
	if err != nil {
		panic(err)
	}

	end, err := strconv.Atoi(char)
	if err != nil {
		panic(err)
	}

	capa := end - beg
	ports := make(chan int, capa)

	go func() {
		for i := beg; i <= end; i++ {
			ports <- i
		}
	}()

	results := make(chan int)

	for i := beg; i <= end; i++ {
		go worker(ports, results)
	}

	for i := beg; i <= end; i++ {
		port := <- results
		if port != 0 {
			openports := append(openports, port)
		}
	}
	close(ports)
	close(results)

	sort.Ints(openports)

	for _, port := range openports {
		fmt.Printf("%d - open\n", port)
	}

}

