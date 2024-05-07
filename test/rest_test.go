package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"rest/internal/config"
	"rest/internal/userProxy"
	"rest/pkg/proto"
	"strings"
	"testing"
)

func TestCreatedUserId(t *testing.T) {
	cfg := config.GetConfig()
	req, err := http.NewRequest("POST", "http://localhost:9091/user", nil)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	bodyString := string(body)
	bodyString = strings.Replace(bodyString, "\n", "", -1)
	var us userProxy.User
	err = json.Unmarshal([]byte(bodyString), &us)
	assert.GreaterOrEqual(t, cfg.User.MinId, us.Id)
	log.Printf("Created user: %v", us.Id)

}
func TestUpdatedUserId(t *testing.T) {
	type tstUser struct {
		Id   int    `json:"Id"`
		Name string `json:"name"`
		Job  string `json:"job"`
	}
	var usr tstUser = tstUser{
		Id:   791,
		Name: "IvanS",
		Job:  "SA",
	}
	usrJ, _ := json.Marshal(usr)
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:9091/user/v3"), bytes.NewReader(usrJ))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	bodyString := string(body)
	bodyString = strings.Replace(bodyString, "\n", "", -1)
	log.Printf("bodyString: %s", bodyString)

	var us tstUser
	err = json.Unmarshal([]byte(bodyString), &us)
	log.Printf("response=%+v", us)
	assert.Equal(t, usr, us)
	log.Printf("Updated: %v", us)

}

func TestGrps(t *testing.T) {
	cfg := config.GetConfig()
	address := fmt.Sprintf("localhost:%s", cfg.Listen.GrpcPort)
	ctx := context.Background()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := proto.NewUserRPCClient(conn)

	personIn, err := client.GetUser(ctx, &proto.GetUserInput{Name: "addsa"})
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	log.Println(personIn)
}
