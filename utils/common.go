/*
 * @Description:
 * @Author: Zjy
 * @Date: 2023-08-23 11:10:17
 * @LastEditTime: 2023-08-23 22:26:23
 * @version: 1.0
 */
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Realurl(s_url string) string {
	r_url := "http://" + URLIP + ":" + PORT + "/data/" + s_url
	return r_url
}

func I64ToTime(num int64) time.Time {
	// 提取时间戳的秒和微秒部分
	seconds := num / 1000
	microseconds := (num % 1000) * 1000
	// 将 Unix 时间戳转换为 time.Time 类型
	timeObj := time.Unix(seconds, microseconds)

	// 格式化为指定的日期时间格式
	//formattedTime := timeObj.Format("2006-01-02 15:04:05.000000")

	return timeObj
}

func TimeToI64(t time.Time) int64 {
	// 将 time.Time 转换为 Unix 时间戳（秒数）
	seconds := t.Unix()

	// 将秒数转换为毫秒级别的时间戳
	milliseconds := seconds*1000 + int64(t.Nanosecond())/int64(time.Millisecond)

	return milliseconds
}

// 返回以时间戳命名的 视频名, 路径, 数据库中的路径
func GetVideoNameAndPath() (name string, path string, DBpath string) {
	// 获取当前时间戳（Unix 时间戳，以纳秒为单位）
	currentTime := time.Now().UnixNano()

	// 将时间戳转换为字符串
	timestampStr := fmt.Sprintf("%d", currentTime)

	name = timestampStr + ".mp4"

	projectBase, _ := os.Getwd()

	path = filepath.Join(projectBase, "data", "videos", name)

	DBpath = filepath.Join("videos", name)

	fmt.Println(path)

	return
}

// 返回视频封面存储名和db名
func GetVideoCoverName(name string) (coverPath string, dbCover string) {
	coverPath = strings.TrimSuffix(name, ".mp4")
	coverPath = coverPath + ".png"
	coverPath = strings.Replace(coverPath, "videos", "bgs", 1)
	index := strings.Index(coverPath, "bgs")
	dbCover = coverPath[index:]
	return
}

// 从字节切片存储视频，弃用
//func SaveVideo(data []byte, savePath string) error {
//	// 打开文件以进行写入
//	file, err := os.Create(savePath)
//	if err != nil {
//		fmt.Println("Error creating file:", err)
//		return ErrSaveVideoFaile
//	}
//	defer file.Close()
//
//	// 将视频数据写入文件
//	_, err = file.Write(data)
//	if err != nil {
//		return ErrSaveVideoFaile
//	}
//	return nil
//}

// 从视频文件中提取指定帧作为封面图，并保存为图片文件
