package functionTest

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"miniDouyin/biz/model/miniDouyin/api"
	"net/http"
	"os"
	"testing"
)

func TestFeed(t *testing.T) {
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

	file, _ := os.Create("Feed_Output.json")
	defer file.Close()

	if err != nil {
		return
	}
	json.Unmarshal(body, resp)

	assert.Equal(t, int32(0), resp.StatusCode)

	if resp.StatusMsg != nil {
		t.Log("返回状态消息为: ", *resp.StatusMsg)
	}

	for _, video := range resp.VideoList {
		t.Log("收到视频: ", video.Title)
	}

	file.Write(body)
}

func Test_Register(t *testing.T) {

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

	file, _ := os.Create("Register_Output.json")
	defer file.Close()
	file.Write(body)
}
