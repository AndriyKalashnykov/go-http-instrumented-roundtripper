package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"
)

type CustomTransport struct {
	http.RoundTripper
	// ... private fields
}

func NewCustomTransport(upstream *http.Transport) *CustomTransport {
	upstream.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// ... other customizations for transport
	return &CustomTransport{upstream}
}

func (ct *CustomTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("Secret", "Blah blah blah")
	// ... customizations for each request

	for i := 1; i <= 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		//reqT := req.WithContext(ctx)
		resp, err = ct.RoundTripper.RoundTrip(req.WithContext(ctx))
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("#%d got timeout will retry - %v", i, err)
			//time.Sleep(time.Duration(100*i) * time.Millisecond)
			continue
		} else {
			break
		}
	}
	//for i := 1; i <= 5; i++ {
	//	resp, err = ct.RoundTripper.RoundTrip(req)
	//	if errors.Is(err, context.DeadlineExceeded) {
	//		fmt.Printf("#%d got timeout will retry - %v\n", i, err)
	//		//time.Sleep(time.Duration(100*i) * time.Millisecond)
	//		continue
	//	} else {
	//		break
	//	}
	//}

	return resp, err
}

func main() {
	transport := NewCustomTransport(http.DefaultTransport.(*http.Transport))
	client := &http.Client{
		Timeout:   8 * time.Second,
		Transport: transport,
	}

	//apiUrl := "https://httpbin.org/delay/10"
	apiUrl := "http://localhost:1338/hello"

	fmt.Printf("Get/begin %q\n", apiUrl)
	start := time.Now()
	resp, err := client.Get(apiUrl)
	if err != nil {
		fmt.Printf("client got error: %v\n", err)
	} else {
		defer resp.Body.Close()
	}
	fmt.Printf("Get/end %q, time cost: %v\n", apiUrl, time.Since(start))

	if resp != nil {
		data, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println("fail to dump resp: %v\n", err)
		}
		fmt.Println(string(data))
	}
}
