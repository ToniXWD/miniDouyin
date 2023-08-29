package rdb

import (
	"context"
	"fmt"
)

// 通过token获取缓存中的记录
func GetVideoById(ID string) (map[string]string, bool) {
	ctx := context.Background()

	// 使用 Exists 方法判断键是否存在
	exists, err := Rdb.Exists(ctx, ID).Result()
	if err != nil || exists != 1 {
		fmt.Println("Error:", err)
		return nil, false
	}

	user, err := Rdb.HGetAll(ctx, ID).Result()
	if err != nil {
		return nil, false
	}
	return user, true
}

//func GetFeedIds(timestamp int64) ([]string, error) {
//	ctx := context.Background()
//
//	stamp := time.UnixMilli(timestamp)
//	cmp_time := stamp.Format(utils.TimeFormat)
//
//	first, err := Rdb.LIndex(ctx, "feed", 0).Result()
//	if err != nil {
//		return nil, err
//	}
//	id_time := strings.Split(first, "_")
//
//	v_time := id_time[0]
//
//	if cmp_time < v_time {
//		//如果list中最早的记录都比时间时间更晚，返回
//		return nil, err
//	}
//
//	// 获取 LIST 的长度
//	length, err := Rdb.LLen(ctx, "feed").Result()
//	if err != nil {
//		fmt.Println("Error:", err)
//		return nil, err
//	}
//
//	// 二分超找 LIST 查找第一个符合条件的元素及其位置
//	var end int64 = length - 1
//	var start int64 = 0
//	for start < end && end < length {
//		mid := start + (end-start)/2
//		element, _ := Rdb.LIndex(ctx, "feed", mid).Result()
//
//		cur_time := strings.Split(element, "_")[0]
//
//		// 检查条件
//		if cur_time > cmp_time {
//			end = mid - 1
//		} else if cur_time < cmp_time {
//			start = mid + 1
//		} else {
//			end = mid
//			break
//		}
//	}
//	start = end - 30
//	if start < 0 {
//		start = 0
//	}
//	vlists, err := Rdb.LRange(ctx, "feed", start, end-1).Result()
//	if err != nil {
//		fmt.Println("Error:", err)
//		return nil, err
//	}
//
//	for idx, item := range vlists {
//		id_time = strings.Split(item, "_")
//
//		v_id := id_time[1]
//		vlists[idx] = v_id
//	}
//	return vlists, nil
//}
//
//func VMap2ApiVidio(client map[string]string, vMap map[string]string) (*api.Video, bool) {
//	tmp, _ := strconv.Atoi(vMap["ID"])
//	ID := int64(tmp)
//	tmp, _ = strconv.Atoi(vMap["Author"])
//	Author := int64(tmp)
//	tmp, _ = strconv.Atoi(vMap["FavoriteCount"])
//	FavoriteCount := int64(tmp)
//	tmp, _ = strconv.Atoi(vMap["CommentCount"])
//	CommentCount := int64(tmp)
//
//	// 查询用户缓存
//	authorMap, find := GetUserById(Author)
//	if !find {
//		return nil, false
//	}
//
//	// 查询是否关注
//	follow, err := IsFollow(client["Token"], ID)
//	if err != nil {
//		return nil, false
//	}
//
//	author := GetApiUserFromMap(authorMap)
//	author.IsFollow = follow
//
//	apiv := &api.Video{
//		ID:            ID,
//		Author:        author,
//		PlayURL:       vMap["PlayURL"],
//		CoverURL:      vMap["CoverURL"],
//		FavoriteCount: FavoriteCount,
//		CommentCount:  CommentCount,
//		Title:         vMap["Title"],
//	}
//
//	return apiv, true
//}
