package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func proxy(defaultUrl *url.URL, altUrl *url.URL, forwardCodes map[int]struct{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		outReq := r.Clone(ctx)

		outReq.URL.Scheme = defaultUrl.Scheme

		// pre-emptively copy request body, in case it needs to be read again
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		outReq.Body = io.NopCloser(bytes.NewReader(body))

		newURL := *r.URL
		newURL.Host = defaultUrl.Host
		newURL.Scheme = defaultUrl.Scheme
		outReq.URL = &newURL
		outReq.Close = false

		res, err := http.DefaultTransport.RoundTrip(outReq)

		var forward bool
		var errStr string
		if err != nil {
			forward = true
			errStr = err.Error()
		} else {
			_, forward = forwardCodes[res.StatusCode]
			errStr = strconv.Itoa(res.StatusCode)
		}

		if forward {
			outReq := r.Clone(ctx)
			outReq.Body = io.NopCloser(bytes.NewReader(body))
			log.Printf("ERROR %s | FORWARD: %s %s", errStr, r.Method, r.URL)
			newURL = *r.URL
			newURL.Host = altUrl.Host
			newURL.Scheme = altUrl.Scheme
			outReq.URL = &newURL
			res, err = http.DefaultTransport.RoundTrip(outReq)
			if err == nil {
				log.Printf("FORWARD: %s %s %s", res.Status, r.Method, r.URL)
			}
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
	}
}
