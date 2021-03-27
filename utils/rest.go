package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	Client = http.Client{}
)

func HttpGet(url string, headers map[string]string, isArr bool) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", url, strings.NewReader(""))
	for k, e := range headers {
		req.Header.Add(k, e)
	}
	res, err := Client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(res.Body)
	var tmp []byte
	if isArr {
		tmp = append(tmp, []byte("{\"guilds\":")...)
		tmp = append(tmp, body...)
		tmp = append(tmp, []byte("}")...)
	} else {
		tmp = append(tmp, body...)
	}

	var data map[string]interface{}
	err = json.Unmarshal(tmp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
