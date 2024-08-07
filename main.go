package main

import (
	"fmt"
	"io"
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
		fmt.Printf("invalid port %s", os.Args[1])
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

	forwardCodes := make(map[int]bool)
	for _, code := range os.Args[4:] {
		i, err := strconv.Atoi(code)
		if err != nil {
			fmt.Printf("invalid code: %s\n", code)
			panic(err)
		}
		forwardCodes[i] = true
	}

	log.Printf("default upstream: %v", defaultUrl)
	log.Printf("alternative upstream: %v", altUrl)
	var s strings.Builder
	s.WriteString("forward response codes: ")
	for code := range forwardCodes {
		s.WriteString(fmt.Sprintf("%d ", code))
	}
	log.Print(s.String())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		outReq := r.Clone(ctx)

		log.Printf("%s %s", r.Method, r.URL)
		outReq.URL.Scheme = defaultUrl.Scheme
		newReq := defaultUrl
		newReq.Path = r.URL.Path
		outReq.URL = newReq
		outReq.Close = false

		res, err := http.DefaultTransport.RoundTrip(outReq)
		if err != nil {
			log.Printf("err: %v", err)
			return
		}

		_, ok := forwardCodes[res.StatusCode]
		if ok {
			log.Printf("ERROR %d | FORWARD: %s %s", res.StatusCode, r.Method, r.URL)
			newReq = altUrl
			newReq.Path = r.URL.Path
			outReq.URL = newReq
			res, err = http.DefaultTransport.RoundTrip(outReq)
		} else {
			log.Printf("%s %s %s", res.Status, r.Method, r.URL)
		}
		if err != nil {
			log.Printf("err: %v", err)
			return
		}

		for key, vals := range res.Header {
			for _, val := range vals {
				w.Header().Add(key, val)
			}
		}

		w.WriteHeader(res.StatusCode)

		_, err = io.Copy(w, res.Body)
		if err != nil {
			return
		}
		res.Body.Close()
	})

	listenStr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(listenStr, nil)
}
