package utils

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func ExtractCover(v_path string, savePath string) {
	// 构建 ffmpeg 命令
	cmd := exec.Command("ffmpeg", "-i", v_path, "-ss", "00:00:01", "-vframes", "1", savePath)

	// 执行命令并捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Infoln("Error running ffmpeg command:", err)
		log.Infoln("Output:", string(output))
		return
	}
	log.Infoln("Image extracted and saved successfully!")
}
