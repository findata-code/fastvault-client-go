package fastvault_client_go

import (
	"fmt"
	"net/http"
	"testing"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "error", 400)
		return
	}

	if r.URL.Path != "/secret" {
		http.Error(w, "error", 400)
		return
	}

	if r.Header.Get("X-Application-Token") == "token" {
		fmt.Fprintf(w, "secret")
		return
	}

	if r.Header.Get("X-Application-Token") == "json" {
		fmt.Fprintf(w, `{"text":"secret"}`)
		return
	}

	http.Error(w, "error", 400)
}

func PostHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST" {
		http.Error(w, "error", 400)
		return
	}

	if r.URL.Path != "/secret" {
		http.Error(w, "error", 400)
		return
	}

	fmt.Fprintf(w, `{"token":"this_is_token"}`)
}

func TestNew(t *testing.T) {
	const url = "http://127.0.0.1"
	client := New(url)
	if client.url != url {
		t.Error("expect", url, "actual", client.url)
	}
}

func TestFastVaultClient_Create(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/secret", PostHandler)

	go http.ListenAndServe(":9801", mux)

	const url = "http://127.0.0.1:9801"

	t.Run("it should return token when call fastvault with body", func(t *testing.T) {
		client := New(url)
		token, err := client.Create("Hello World")
		if err != nil {
			t.Error(err)
		}

		if token != "this_is_token" {
			t.Error("expect this_is_token actual", token)
		}
	})
}

func TestFastVaultClient_GetString(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/secret", GetHandler)
	go http.ListenAndServe(":9800", mux)

	const url = "http://127.0.0.1:9800"

	t.Run("it should return secret when call to fastvault", func(t *testing.T) {
		client := New(url)
		res, err := client.GetString("token")
		if err != nil {
			t.Error(err)
		}

		if res != "secret" {
			t.Error("expect secret actual", res)
		}
	})

	t.Run("it should fulfill struct when call to fastvault", func(t *testing.T) {
		type helloWorld struct {
			Text string `json:"text"`
		}

		var v helloWorld

		client := New(url)

		err := client.GetJson("json", &v)
		if err != nil {
			t.Error(err)
		}

		if v.Text != "secret" {
			t.Error(err)
		}
	})
}
