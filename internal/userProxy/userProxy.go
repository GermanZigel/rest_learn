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
	Name    string `json:"name"`
	Job     string `json:"job"`
	Created string `json:"Created"`
	Id      int    `json:"Id"`
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

	myUser.Id = int(gjson.Get(string(body), "id").Int())
	log.Println("Создан пользователь с ID =", myUser.Id)
	myUser.Name = gjson.Get(string(body), "name").Str
	myUser.Job = gjson.Get(string(body), "job").Str
	myUser.Created = gjson.Get(string(body), "createdAt").Str
	myUser.Comment = "User created sucess"
	log.Println("начало проверки пользователя с ID =", myUser.Id)
	cfg := config.GetConfig()
	for myUser.Id < cfg.Listen.MitId {
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
