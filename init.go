package nsqTask

import (
	"github.com/nsqio/go-nsq"
	"sync"
)

var (
	initOnce  sync.Once
	Producers map[string]*nsq.Producer
)

const TestTask1 = "testTask1"
const TestTask2 = "testTask2"

func SetUp() {
	initOnce.Do(func() {
		Producers = make(map[string]*nsq.Producer, 2)
		sendTask(TaskConfig{
			TopicName:        TestTask1,
			ConsumerTotal:    3,
			NsqAddress:       "127.0.0.1:4150",
			NsqLookupAddress: "127.0.0.1:4161",
			ChannelName:      "c1",
			Handler:          &Task1Handler{},
		})
		sendTask(TaskConfig{
			TopicName:        TestTask2,
			ConsumerTotal:    3,
			NsqAddress:       "127.0.0.1:4150",
			NsqLookupAddress: "127.0.0.1:4161",
			ChannelName:      "c1",
			Handler:          &Task2Handler{},
		})
	})
}
