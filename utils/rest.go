package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	client http.Client
)

func InitHttp() {
	client = http.Client{}
}

func HttpGet(url string, headers map[string]string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", url, strings.NewReader(""))
	for k, e := range headers {
		req.Header.Add(k, e)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(res.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)
	return data, nil
}
