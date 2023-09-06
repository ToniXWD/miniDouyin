package functionTest

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_FavoriteAction(t *testing.T) {

	url := address + "/douyin/favorite/action/?token=test2123456&video_id=4&action_type=2"
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
	CreateLogFile("FavoriteAction_2", body)
}

func Test_FavoriteList(t *testing.T) {

	url := address + "/douyin/favorite/list/?user_id=4&token=test2123456"
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

	CreateLogFile("FavoriteList", body)
}

// TODO：未测试
func Test_Comment(t *testing.T) {

	url := address + "/douyin/comment/action/?token=test2123456&video_id=4&action_type=1"
	method := "POST"

	client := &http.Client{}
	message := "你好啊，测试一下评论"
	reqbody := bytes.NewReader([]byte(message))
	req, err := http.NewRequest(method, url, reqbody)

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
	CreateLogFile("Comment", body)
}

func Test_CommentList(t *testing.T) {
	url := address + "/douyin/comment/list/?token=test2123456&video_id=4"
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	CreateLogFile("CommentList", body)
}
