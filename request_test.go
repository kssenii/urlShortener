package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestValidateURL(t *testing.T) {

	url := "http://localhost:9090/encode"

	values := map[string]string{"url": "http://www.vkonl.com"}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if string(resp.Status) == "200 OK" {
		t.Error("Exited without error, but validation must fail")
	}
}

func TestGenerateURL(t *testing.T) {

	url := "http://localhost:9090/encode"

	values := map[string]string{"url": "http://www.instagram.com"}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if string(resp.Status) != "200 OK" {
		t.Error("Exited with ", string(resp.Status))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	t.Logf("response Body: %s", string(body))
}

func TestCustomURL(t *testing.T) {

	url := "http://localhost:9090/encode"

	values := map[string]string{"url": "http://www.yandex.ru", "short": "ya"}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if string(resp.Status) != "200 OK" {
		t.Error("Exited with ", string(resp.Status))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	ref := `{"url":"http://localhost:9090/ya"}`

	if strings.Trim(string(body), " \r\n") != ref {
		t.Error("Expected:", ref, "Got:", string(body))
	}
	t.Logf("response Body: %s", string(body))
}

func TestValidateShortURL(t *testing.T) {

	url := "http://localhost:9090/decode"

	values := map[string]string{"url": "http://localhost:9/vk"}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if string(resp.Status) == "200 OK" {
		t.Error("Exited without error, but validation must fail")
	}
}

func TestDecodeURL(t *testing.T) {

	url := "http://localhost:9090/decode"

	values := map[string]string{"url": "http://localhost:9090/ya"}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if string(resp.Status) != "200 OK" {
		t.Error("Exited with ", string(resp.Status))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	ref := `{"url":"http://www.yandex.ru"}`

	if strings.Trim(string(body), " \r\n") != ref {
		t.Error("Expected:", ref, "Got:", string(body))
	}
	t.Logf("response Body: %s", string(body))
}

func TestRedirectURL(t *testing.T) {

	url := "http://localhost:9090/ya"

	req, err := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if string(resp.Status) != "200 OK" {
		t.Error("Exited with ", string(resp.Status))
	}
}
