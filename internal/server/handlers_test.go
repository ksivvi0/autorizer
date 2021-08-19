package server

import (
	"authorizer/internal/services"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

type headers struct {
	key   string
	value string
}

type reqParams struct {
	url          string
	method       string
	responseCode int
	authHeader   *headers
	body         map[string]string
}

type tokenPair struct {
	AccessToken    string    `json:"access_token,omitempty"`
	AccessUID      string    `json:"access_token_uid"`
	AccessExpired  time.Time `json:"access_token_expired,omitempty"`
	RefreshToken   string    `json:"refresh_token"`
	RefreshUID     string    `json:"refresh_token_uid"`
	RefreshExpired time.Time `json:"refresh_token_expired,omitempty"`
}

var (
	bearer string
	srv    *Server
	tPair  *tokenPair = new(tokenPair)
)

func init() {
	_ = os.Setenv("MONGO_URI", "mongodb+srv://tester:tester123@authorizer.ik7zm.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	_ = os.Setenv("SERVER_ADDR", "localhost:9000")
	_ = os.Setenv("LOG_PATH", "./test.log")
	_ = os.Setenv("CRYPTO_SECRET", "asasdl")
	_ = os.Setenv("MONGO_DB", "authorizer")
	_ = os.Setenv("MONGO_COLLECTION", "main")
	_ = os.Setenv("MAGIC_WORD", "test")

	l, err := services.NewLoggerInstance(os.Getenv("LOG_PATH"))
	if err != nil {
		panic(err)
	}
	a := services.NewAuthInstance()
	s, err := services.NewStoreInstance(os.Getenv("MONGO_URI"))
	if err != nil {
		panic(err)
	}
	if err = s.Ping(); err != nil {
		panic(err)
	}

	svcs := services.NewServices(l, a, s)
	srv, err = NewServerInstance(os.Getenv("SERVER_ADDR"), svcs, false)
	if err != nil {
		panic(err)
	}
	go func() {
		if err := srv.Run(); err != nil {
			panic(err)
		}
	}()
}
func TestCreateTokenHandler(t *testing.T) {
	cases := []reqParams{
		{url: "/auth/tokens", method: "POST", body: map[string]string{"magic_word": "test"}, responseCode: 200},
		{url: "/auth/tokens", method: "POST", body: map[string]string{"magic_word": "bad_word"}, responseCode: 403},
	}
	httpClient := http.Client{}

	for _, v := range cases {
		var body io.Reader

		if v.body != nil {
			rawBody, _ := json.Marshal(v.body)
			body = bytes.NewBuffer(rawBody)
		}

		req, err := http.NewRequest(v.method, fmt.Sprintf("http://%s%s", os.Getenv("SERVER_ADDR"), v.url), body)
		if err != nil {
			t.Error(err)
		}

		if v.authHeader != nil {
			req.Header.Add(v.authHeader.key, v.authHeader.value)
		}
		response, err := httpClient.Do(req)
		if err != nil {
			t.Errorf("URL: %s, %v", v.url, err)
		}

		if response.StatusCode == 200 {
			if err = json.NewDecoder(response.Body).Decode(tPair); err != nil {
				t.Error(err)
			}
			bearer = fmt.Sprintf("Bearer %s", tPair.AccessToken)
		}
		if response.StatusCode != v.responseCode {
			t.Errorf("URL: %s, receive code: %d, wants: %d", v.url, response.StatusCode, v.responseCode)
		}
	}
}

func TestPingHandler(t *testing.T) {
	cases := []reqParams{
		{url: "/api/ping", method: "GET", responseCode: 401},
		{url: "/api/ping", method: "POST", responseCode: 404},
		{url: "/api/ping", method: "GET", authHeader: &headers{"Authorization", bearer}, responseCode: 200},
	}
	httpClient := http.Client{}

	for _, v := range cases {
		var body io.Reader

		if v.body != nil {
			rawBody, _ := json.Marshal(v.body)
			body = bytes.NewBuffer(rawBody)
		}

		req, err := http.NewRequest(v.method, fmt.Sprintf("http://%s%s", os.Getenv("SERVER_ADDR"), v.url), body)
		if err != nil {
			t.Error(err)
		}

		if v.authHeader != nil {
			req.Header.Add(v.authHeader.key, v.authHeader.value)
		}
		response, err := httpClient.Do(req)
		if err != nil {
			t.Errorf("URL: %s, %v", v.url, err)
		}

		if response.StatusCode != v.responseCode {
			t.Errorf("URL: %s with headers %v, receive code: %d, wants: %d", v.url, v.authHeader.key, response.StatusCode, v.responseCode)
		}
	}
}
