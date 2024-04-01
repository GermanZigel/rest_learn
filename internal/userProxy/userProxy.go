package userProxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"rest/internal/config"
)

type User struct {
	Id      int    `json:"Id"`
	Name    string `json:"name"`
	Job     string `json:"job"`
	Created string `json:"Created,omitempty"`
	Comment string `json:"Comment,omitempty"`
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
	cfg := config.GetConfig()
	myUser := &User{}
	userName, userJob := cfg.User.Name, cfg.User.Job
	myUser.Set(userName, userJob)

	js, err := jsonStr(myUser)
	if err != nil {
		fmt.Println("Ошибка при получении JSON строки:", err)
		return nil
	}

	httpPostUrl := cfg.Listen.HOST
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

	myUser.Id = int(gjson.Get(string(body), "id").Int())
	log.Println("Создан пользователь с ID =", myUser.Id)
	myUser.Name = gjson.Get(string(body), "name").Str
	myUser.Job = gjson.Get(string(body), "job").Str
	myUser.Created = gjson.Get(string(body), "createdAt").Str
	myUser.Comment = "User created sucess"
	log.Println("начало проверки пользователя с ID =", myUser.Id)
	for myUser.Id < cfg.User.MinId {
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		myUser.Id = int(gjson.Get(string(body), "id").Int())
		myUser.Name = gjson.Get(string(body), "name").Str
		myUser.Job = gjson.Get(string(body), "job").Str
		myUser.Created = gjson.Get(string(body), "createdAt").Str
		myUser.Comment = "User created success"
		log.Println(myUser)

		// Create a new request for the next iteration
		req, err = http.NewRequest("POST", httpPostUrl, bytes.NewBuffer(js))
		req.Header.Set("Content-type", "application/json; charset=UTF-8")
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	return myUser
}
