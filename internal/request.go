package internal

import (
	"io/ioutil"
	"net/http"
)

func TryRequest(path string) (string, error) {
	client := http.Client{}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body)[:100], nil
}
