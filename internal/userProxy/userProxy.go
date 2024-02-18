package userProxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

type User struct {
	Name    string `json:"name"`
	Job     string `json:"job"`
	Created string `json:"Created"`
	Id      string `json:"Id"`
	Comment string `json:"Comment"`
}

func jsonStr(m *User) ([]byte, error) {
	xx, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Ошибка при маршалинге JSON:", err)
		return nil, err
	}
	return xx, nil
}

func (b *User) Set(Name, Job string) {
	b.Name = Name
	b.Job = Job
}

func Setter() *User {
	myUser := &User{}
	myUser.Set("Ht", "Worker")

	js, err := jsonStr(myUser)
	if err != nil {
		fmt.Println("Ошибка при получении JSON строки:", err)
		return nil
	}

	httpPostUrl := "https://reqres.in/api/users"
	req, err := http.NewRequest("POST", httpPostUrl, bytes.NewBuffer(js))
	req.Header.Set("Content-type", "application/json; charset=UTF-8")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	myUser.Id = gjson.Get(string(body), "id").Str
	myUser.Name = gjson.Get(string(body), "name").Str
	myUser.Job = gjson.Get(string(body), "job").Str
	myUser.Created = gjson.Get(string(body), "createdAt").Str
	myUser.Comment = "User created sucess"
	return myUser

	return myUser

}

/*
{"name":"Ht","job":"Worker","id":"102","createdAt":"2024-02-17T15:59:13.369Z"}

*/
