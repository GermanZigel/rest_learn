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
	address := fmt.Sprintf("http://localhost:%s%s", cfg.Listen.HttpPort, cfg.Listen.URI_Once)
	req, err := http.NewRequest("POST", address, nil)
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
	assert.GreaterOrEqual(t, us.Id, cfg.User.MinId)
	log.Printf("Created user: %v", us.Id)

}
func TestUpdatedUserId(t *testing.T) {
	cfg := config.GetConfig()
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
	address := fmt.Sprintf("http://localhost:%s%s", cfg.Listen.HttpPort, cfg.Listen.URI_Once)
	req, err := http.NewRequest("PUT", address, bytes.NewReader(usrJ))
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

	personIn, err := client.GetUser(ctx, &proto.GetUserInput{Id: 791})
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	log.Println(personIn)
	assert.Less(t, cfg.User.MinId, int(personIn.Id))
	assert.Equal(t, cfg.User.Job, personIn.Job)
}
