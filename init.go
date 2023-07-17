package nsqQueue

import (
	"github.com/nsqio/go-nsq"
	"sync"
)

var (
	initOnce   sync.Once
	Producers  map[string]*nsq.Producer
	TaskLogger sync.Map
)

const CalculateMardTopic = "calculateMardTask"
const UploadFileToQiniuTopic = "uploadFileToQiniuTask"

func SetUp() {
	initOnce.Do(func() {
		setLogger()
		Producers = make(map[string]*nsq.Producer, 2)
		sendTask(TaskConfig{
			TopicName:        CalculateMardTopic,
			ConsumerTotal:    3,
			NsqAddress:       "127.0.0.1:4150",
			NsqLookupAddress: "127.0.0.1:4161",
			ChannelName:      "c1",
			Handler:          &CalculateMardHandler{},
		})
		sendTask(TaskConfig{
			TopicName:        UploadFileToQiniuTopic,
			ConsumerTotal:    3,
			NsqAddress:       "127.0.0.1:4150",
			NsqLookupAddress: "127.0.0.1:4161",
			ChannelName:      "c1",
			Handler:          &UploadFileToQiniuHandler{},
		})
	})
}
