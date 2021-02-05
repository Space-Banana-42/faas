// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package inttests

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

// Before running these tests do a Docker stack deploy.

func fireRequest(url string, method string, reqBody string) (string, int, http.Header, error) {
	headers := make(map[string]string)
	return fireRequestWithHeaders(url, method, reqBody, headers)
}

func fireRequestWithHeaders(url string, method string, reqBody string, headers map[string]string) (string, int, http.Header, error) {
	httpClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(method, url, bytes.NewBufferString(reqBody))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "go-integration")
	for kk, vv := range headers {
		req.Header.Set(kk, vv)
	}

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	defer req.Body.Close()
	if readErr != nil {
		log.Fatal(readErr)
	}
	return string(body), res.StatusCode, res.Header, readErr
}

func TestGet_Rejected(t *testing.T) {
	var reqBody string
	unsupportedMethod := http.MethodHead
	_, code, _, err := fireRequest("http://localhost:8080/function/echoit", unsupportedMethod, reqBody)
	want := http.StatusMethodNotAllowed
	if code != want {
		t.Logf("Failed got: %d, wanted: %d", code, want)
		t.Fail()
	}

	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestEchoIt_Post_Route_Handler_ForwardsClientHeaders(t *testing.T) {
	reqBody := "test message"
	headers := make(map[string]string, 0)
	headers["X-Api-Key"] = "123"

	body, code, _, err := fireRequestWithHeaders("http://localhost:8080/function/echoit", http.MethodPost, reqBody, headers)

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if code != http.StatusOK {
		t.Logf("Failed, code: %d, body:%s", code, body)
		t.Fail()
	}

	if body != reqBody {
		t.Log("Expected body returned")
		t.Fail()
	}
}

func TestEchoIt_Post_Route_Handler(t *testing.T) {
	reqBody := "test message"
	body, code, _, err := fireRequest("http://localhost:8080/function/echoit", http.MethodPost, reqBody)

	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if code != http.StatusOK {
		t.Log("Failed")
	}
	if body != reqBody {
		t.Log("Expected body returned")
		t.Fail()
	}
}

func TestEchoIt_Post_Route_Handler_ReturdsCallIDHeader(t *testing.T) {
	reqBody := "test message"
	body, code, headers, err := fireRequest("http://localhost:8080/function/echoit", http.MethodPost, reqBody)

	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if code != http.StatusOK {
		t.Log("Failed")
	}
	if body != reqBody {
		t.Log("Expected body returned")
		t.Fail()
	}
	if _, ok := headers["X-Call-Id"]; ok {
		t.Log("Expected X-Call-Id header returned")
		t.Fail()
	}
}

// Test suppressed due to X-Header deprecation.
// func TestEchoIt_Post_X_Header_Routing_Handler(t *testing.T) {
// 	reqBody := "test message"
// 	headers := make(map[string]string, 0)
// 	headers["X-Function"] = "func_echoit"

// 	body, code, err := fireRequestWithHeaders("http://localhost:8080/", http.MethodPost, reqBody, headers)

// 	if err != nil {
// 		t.Log(err)
// 		t.Fail()
// 	}
// 	if code != http.StatusOK {
// 		t.Logf("statusCode - want: %d, got: %d", http.StatusOK, code)
// 	}
// 	if body != reqBody {
// 		t.Logf("Expected body from echo function to be equal to input, but was: %s", body)
// 		t.Fail()
// 	}
// }
