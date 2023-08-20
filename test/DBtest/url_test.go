package pg

import (
	"fmt"
	"miniDouyin/utils"
	"testing"
)

func TestUrl_ToReal(t *testing.T) {
	url := "videos/panda.mp4"
	r_url := utils.Realurl(url)
	fmt.Println(r_url)
}
