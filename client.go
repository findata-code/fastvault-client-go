package fastvault_client_go

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const TOKEN = "X-Application-Token"

type FastVaultClient struct {
	url string
	httpClient http.Client
}

func New(url string) *FastVaultClient {
	return &FastVaultClient{
		url,
		http.Client{},
	}
}

func (c *FastVaultClient) GetString(token string) (string, error) {
	b, err := c.get(token)
	return string(b), err
}

func (c *FastVaultClient) GetByte(token string) ([]byte, error) {
	return c.get(token)
}

func (c *FastVaultClient) GetJson(token string, v interface{}) error {
	b, err := c.get(token)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

func (c *FastVaultClient) get(token string) ([]byte, error) {
	url := strings.TrimSuffix(c.url, "/")
	url = fmt.Sprintf("%s/%s", url, "secret")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(TOKEN, token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}
