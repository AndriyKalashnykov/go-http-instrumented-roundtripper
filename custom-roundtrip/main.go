package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
)

type transport struct {
	current *http.Request
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.current = req

	//switch req.URL.String() {
	//case "https://github.com":
	//	responseBody := "This is github.com stub"
	//	respReader := io.NopCloser(strings.NewReader(responseBody))
	//	resp := http.Response{
	//		StatusCode:    http.StatusOK,
	//		Body:          respReader,
	//		ContentLength: int64(len(responseBody)),
	//		Header: map[string][]string{
	//			"Content-Type": {"text/plain"},
	//		},
	//	}
	//	return &resp, nil
	//
	//case "https://example.com":
	//	return http.DefaultTransport.RoundTrip(r)
	//
	//default:
	//	return nil, errors.New("Request URL not supported by stub")
	//}

	return http.DefaultTransport.RoundTrip(req)
}

func main() {

	url := "https://gophercoding.com"
	fmt.Printf("Connecting to: %s\n", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	clientTrace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			log.Printf("Connection was reused: %t", info.Reused)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			log.Printf("DNS Info: %+v\n", dnsInfo)
		},
	}
	// Call once, new connection
	MakeHttpCall(clientTrace, req)

	// Call again, should reuse the connection
	MakeHttpCall(clientTrace, req)
}

// MakeHttpCall is an example of making a http request, while logging any DNS info
// received and if the connection was established afresh, or re-used.
func MakeHttpCall(clientTrace *httptrace.ClientTrace, req *http.Request) {
	t := &transport{}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), clientTrace))
	//_, err := http.DefaultTransport.RoundTrip(req)
	//return err
	client := &http.Client{
		Transport: t,
	}
	if _, err := client.Do(req); err != nil {
		log.Fatal(err)
	}
}
