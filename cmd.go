package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		panic("cmd port defaultUrl altUrl fwdCodes...")
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("invalid port %s\n", os.Args[1])
		panic(err)
	}

	defaultUrl, err := url.Parse(os.Args[2])
	if err != nil {
		fmt.Println("invalid defaultUrl")
		panic(err)
	}
	altUrl, err := url.Parse(os.Args[3])
	if err != nil {
		fmt.Println("invalid altUrl")
		panic(err)
	}

	forwardCodes := make(map[int]struct{})
	for _, code := range os.Args[4:] {
		i, err := strconv.Atoi(code)
		if err != nil {
			fmt.Printf("invalid code: %s\n", code)
			panic(err)
		}
		forwardCodes[i] = struct{}{}
	}

	log.Printf("default upstream: %v", defaultUrl)
	log.Printf("alternative upstream: %v", altUrl)
	var s strings.Builder
	s.WriteString("forward response codes: ")
	for code := range forwardCodes {
		s.WriteString(fmt.Sprintf("%d ", code))
	}
	log.Print(s.String())

	listenStr := fmt.Sprintf(":%d", port)
	http.HandleFunc("/", proxy(defaultUrl, altUrl, forwardCodes))
	http.ListenAndServe(listenStr, nil)
}
