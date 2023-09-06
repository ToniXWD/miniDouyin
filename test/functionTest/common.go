package functionTest

import "os"

const address = "http://172.20.10.4:8889"

func CreateLogFile(name string, body []byte) {
	n := name + "_OutPut.json"
	file, _ := os.Create(n)
	defer file.Close()
	file.Write(body)
}
