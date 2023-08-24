package IOTest

import (
	"github.com/stretchr/testify/assert"
	"miniDouyin/utils"
	"testing"
)

func TestExtract_Cover(t *testing.T) {
	v_path := "../../data/videos/1692607008596902615.mp4"
	coverPath, dbCover := utils.GetVideoCoverName(v_path)
	//fmt.Println("coverPath = ", coverPath)
	//fmt.Println("dbCover = ", dbCover)
	assert.Equal(t, coverPath, "../../data/bgs/1692607008596902615.png")
	assert.Equal(t, dbCover, "bgs/1692607008596902615.png")
	utils.ExtractCover(v_path, coverPath)
}
