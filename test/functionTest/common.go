package functionTest

import (
	"miniDouyin/utils"
	"os"
)

var address = "http://" + utils.URLIP + ":" + utils.PORT

func CreateLogFile(name string, body []byte) {
	n := name + "_OutPut.json"
	file, _ := os.Create(n)
	defer file.Close()
	file.Write(body)
}
