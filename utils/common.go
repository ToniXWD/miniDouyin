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
	"reflect"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var TimeFormat_S string = "2006-01-02 15:04:05"
var TimeFormat_MS string = "2006-01-02 15:04:05.000"
var TimeFormat_MCS string = "2006-01-02 15:04:05.000000"

func Realurl(s_url string) string {
	r_url := "http://" + URLIP + ":" + PORT + "/data/" + s_url
	return r_url
}

func TimeToI64(t time.Time) int64 {
	// 将 time.Time 转换为 Unix 时间戳（秒数）
	seconds := t.Unix()

	// 将秒数转换为毫秒级别的时间戳
	milliseconds := seconds*1000 + int64(t.Nanosecond())/int64(time.Millisecond)

	return milliseconds
}

func Str2Time(fieldValue string) (time.Time, error) {
	s_len := len(TimeFormat_S)
	ms_len := len(TimeFormat_MS)
	mcs_len := len(TimeFormat_MCS)

	time_len := len(fieldValue)
	var t time.Time
	var err error = nil
	if time_len == s_len {
		// 精度为秒
		t, err = time.Parse(TimeFormat_S, fieldValue)
	} else if time_len == ms_len {
		// 精度为毫秒
		t, err = time.Parse(TimeFormat_MS, fieldValue)
	} else if time_len == mcs_len {
		// 精度为微秒
		t, err = time.Parse(TimeFormat_MCS, fieldValue)
	}
	return t, err
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

	log.Infoln("上传文件存储路径为：", path)

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
// func SaveVideo(data []byte, savePath string) error {
//	// 打开文件以进行写入
//	file, err := os.Create(savePath)
//	if err != nil {
//		log.Debugln("Error creating file:", err)
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
// }
// 从视频文件中提取指定帧作为封面图，并保存为图片文件

// 将结构体转化为
func StructToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	value := reflect.ValueOf(data)

	// 确保传入的是结构体指针
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return result
	}

	value = value.Elem()
	typ := value.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == "Deleted" {
			continue
		}
		fieldValue := value.Field(i).Interface()

		if field.Name == "CreatedAt" {
			// 将时间字段转换为指定格式的字符串
			if createdAt, ok := fieldValue.(time.Time); ok {
				result[field.Name] = createdAt.Format("2006-01-02 15:04:05.000")
			} else {
				result[field.Name] = fieldValue
			}
		} else {
			result[field.Name] = fieldValue
		}
	}

	return result
}
