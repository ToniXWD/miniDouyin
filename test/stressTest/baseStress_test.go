package stressTest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"miniDouyin/biz/model/miniDouyin/api"
	"net/http"
	"testing"
)

func SendFeedRequest() {
	url := "http://172.29.172.57:8889/douyin/feed/"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp := new(api.FeedResponse)

	json.Unmarshal(body, resp)

	if resp.StatusCode != 0 {
		panic(1)
	}
}

func SendRegisetr() {
	url := "http://172.29.172.57:8889/douyin/user/register/?username=XXX&password=111111"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(body) == 0 {
		panic(1)
	}

	//file, _ := os.Create("Register_Output.json")
	//defer file.Close()
	//file.Write(body)}
}

func SendGetUserInfo() {

	url := "http://172.29.172.57:8889/douyin/user/?user_id=4&token=toni123456"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(body) == 0 {
		panic(1)
	}
	//fmt.Println(string(body))
}

func Test_base(t *testing.T) {
	for i := 0; i < 100000; i++ {
		//go test_Feed()
		go SendGetUserInfo()
	}
}
