/*
 * @Description:
 * @Author: Zjy
 * @Date: 2023-09-01 19:20:58
 * @LastEditTime: 2023-09-01 19:35:49
 * @version: 1.0
 */
package test

import (
	"context"
	"fmt"
	"miniDouyin/biz/dal/rdb"
	"testing"
)

func TestRedis(t *testing.T) {
	// 测试redis连接
	rdb.Init()

	ctx := context.Background()
	value, _ := rdb.Rdb.Ping(ctx).Result()
	fmt.Println(value)

}
