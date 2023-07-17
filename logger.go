package nsqQueue

import (
	"fmt"
	"log"
	"os"
)

func setLogger() {
	setCalculateMardTaskLogger()
	setUploadFileToQiniuTaskLogger()
}

func setCalculateMardTaskLogger() {
	basePath, err := os.Getwd()
	if err != nil {
		fmt.Printf("get base path failed, err:%v\n", err)
		panic(err)
	}
	baselogPath := basePath + "/cgmlog/nsqQueue"
	file, err := os.OpenFile(baselogPath+"/calculateMardTask.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("打开任务队列日志文件失败,请求后端排除错误")
		panic(err.Error())
	}
	CalculateMardTaskLogger := log.New(file, "[calculateMardTask]", 125)
	TaskLogger.Store(CalculateMardTopic, CalculateMardTaskLogger)
	file.Close()
}

func setUploadFileToQiniuTaskLogger() {
	basePath, err := os.Getwd()
	if err != nil {
		fmt.Printf("get base path failed, err:%v\n", err)
		panic(err)
	}
	baselogPath := basePath + "/cgmlog/nsqQueue"
	file, err := os.OpenFile(baselogPath+"/uploadFileToQiniu.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("打开任务队列日志文件失败,请求后端排除错误")
		panic(err.Error())
	}
	UploadFileToQiniuTaskLogger := log.New(file, "[uploadFileToQiniuTask]", 125)
	TaskLogger.Store(UploadFileToQiniuTopic, UploadFileToQiniuTaskLogger)
	file.Close()
}

func GetTaskLogger(topicName string) *log.Logger {
	if l, ok := TaskLogger.Load(topicName); !ok {
		return nil
	} else {
		return l.(*log.Logger)
	}
}
