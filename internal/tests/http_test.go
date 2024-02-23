package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {
	const wantCode = 200
	httpPostUrl := "http://localhost:1234/users/123"
	req, err := http.NewRequest("POST", httpPostUrl, nil)
	req.Header.Set("Content-type", "application/json; charset=UTF-8")
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

}
