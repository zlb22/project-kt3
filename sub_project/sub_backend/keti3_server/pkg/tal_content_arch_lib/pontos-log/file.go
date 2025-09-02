package plog

import (
	"io"
	"log"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// NewFileWriterWithLumberjack 根据规则写入固定目录
// a.log -> 20060102_a.log
// 轮转 之后 为  20060102T15:04:05.000_a.log
//func NewFileWriterWithLumberjack(slug string) io.Writer {
//	folderPath, fileName := path.Split(slug)
//
//	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
//		err = os.Mkdir(folderPath, 0777)
//		if err != nil {
//			fmt.Println(err)
//			return os.Stderr
//		}
//	}
//
//	file := path.Join(folderPath, fileName)
//	//default using 1 GB + 7 days rotation
//	rotation := &lumberjack.Logger{
//		Filename:  file,
//		MaxSize:   1024, //MB
//		MaxAge:    7,    //days
//		LocalTime: true,
//	}
//	return rotation
//}

func NewFileWriter(slug string) io.Writer {
	folderPath, fileName := path.Split(slug)

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err = os.Mkdir(folderPath, 0777)
		if err != nil {
			log.Println("mkdir for log path error:", err)
			return os.Stderr
		}
	}

	file := path.Join(folderPath, "%Y%m%d-%H_"+fileName)
	rl, err := rotatelogs.New(file, rotatelogs.WithClock(rotatelogs.Local), rotatelogs.WithRotationTime(time.Hour))
	if err != nil {
		log.Println("rotate new log file error:", err)
		return os.Stderr
	}

	return rl
}
