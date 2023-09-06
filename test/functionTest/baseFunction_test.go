package functionTest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"mime/multipart"
	"miniDouyin/biz/model/miniDouyin/api"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestFeed(t *testing.T) {
	url := address + "/douyin/feed/"
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

	url := address + "/douyin/user/register/?username=XXX&password=111111"
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, _ := os.Create("Register_Output.json")
	defer file.Close()
	file.Write(body)
}

func Test_Login(t *testing.T) {

	url := address + "/douyin/user/login/?username=test2&password=123456"
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

	file, _ := os.Create("Login_Output.json")
	defer file.Close()
	file.Write(body)
}

func Test_GetUserInfo(t *testing.T) {

	url := address + "/douyin/user/?user_id=4&token=test2123456"
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
	fmt.Println(string(body))
	file, _ := os.Create("GetUserInfo_OutPut.json")
	defer file.Close()
	file.Write(body)
}

// TODO: 未测试
func Test_Publish(t *testing.T) {

	url := address + "/douyin/publish/action/"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("")
	defer file.Close()
	part1,
		errFile1 := writer.CreateFormFile("data", filepath.Base(""))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return
	}
	_ = writer.WriteField("token", "")
	_ = writer.WriteField("title", "")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func Test_PublishList(t *testing.T) {

	url := address + "/douyin/publish/list/?token=test2123456&user_id=4"
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
	fmt.Println(string(body))
	CreateLogFile("PublishList", body)
}
