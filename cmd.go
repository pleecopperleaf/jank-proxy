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
		fmt.Println("cmd port defaultUrl altUrl fwdCodes...")
		return
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("invalid port %s\n", os.Args[1])
		return
	}

	defaultUrl, err := url.Parse(os.Args[2])
	if err != nil {
		fmt.Println("invalid defaultUrl")
		return
	}
	altUrl, err := url.Parse(os.Args[3])
	if err != nil {
		fmt.Println("invalid altUrl")
		return
	}

	forwardCodes := make(map[int]struct{})
	for _, code := range os.Args[4:] {
		i, err := strconv.Atoi(code)
		if err != nil {
			fmt.Printf("invalid code: %s\n", code)
			return
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
