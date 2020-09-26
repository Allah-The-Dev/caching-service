package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func sendGetRequest(client *http.Client, addr string) {
	res, err := client.Get(addr)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 200 {
		panic("request failed")
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	err = res.Body.Close()
	if err != nil {
		panic(err)
	}
}

func sendPostRequest(client *http.Client, addr string, body map[string]string) {

	jsonValue, _ := json.Marshal(body)

	resp, err := http.Post(addr, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic("request failed")
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = resp.Body.Close()
	if err != nil {
		panic(err)
	}
}

func BenchMarkGetEmployees(b *testing.B) {

	client := &http.Client{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sendGetRequest(client, "http://127.0.0.1:8080/api/v1/employee")
	}
}

func BenchMarkGetEmployee(b *testing.B) {

	client := &http.Client{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sendGetRequest(client, "http://127.0.0.1:8080/api/v1/employee/raja")
	}
}

func BenchMarkPostEmployee(b *testing.B) {

	client := &http.Client{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reqBody := map[string]string{"name": "foo", "unit": "bar"}
		sendPostRequest(client, "http://127.0.0.1:8080/api/v1/employee", reqBody)
	}
}
