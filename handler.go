package nsqTask

import (
	"fmt"
	"github.com/nsqio/go-nsq"
)

type Task1Handler struct {
	Title string
}

func (m *Task1Handler) HandleMessage(msg *nsq.Message) (err error) {
	fmt.Printf("%s recv from %v, msg:%v\n", m.Title, msg.NSQDAddress, string(msg.Body))
	return
}

type Task2Handler struct {
	Title string
}

func (m *Task2Handler) HandleMessage(msg *nsq.Message) (err error) {
	fmt.Printf("%s recv from %v, msg:%v\n", m.Title, msg.NSQDAddress, string(msg.Body))
	return
}
