package main

import (
	"log"
	"nsqTask"
	"os"
)

func main() {
	nsqTask.SetUp() //initialization,Generally at the framework entry file
	producer := nsqTask.GetProducer(nsqTask.TestTask1)
	msg := "msg 1"
	//push data,Generally in business logic
	file, err := os.OpenFile("./error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)
	err = producer.Publish(nsqTask.TestTask1, []byte(msg))
	if err != nil {
		log.Printf("push message failed:%v", err.Error())
	} else {
		log.Printf("push message successful")
	}

	msg = "msg 2"
	err = producer.Publish(nsqTask.TestTask2, []byte(msg))
	if err != nil {
		log.Printf("push message failed:%v", err.Error())
	} else {
		log.Printf("push message successful")
	}
}
