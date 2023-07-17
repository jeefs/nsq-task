package nsqTask

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type TaskConfig struct {
	ConsumerTotal    int
	NsqAddress       string
	NsqLookupAddress string
	TopicName        string
	ChannelName      string
	Handler          nsq.Handler
}

// register task
func sendTask(config TaskConfig) {
	producer, err := initProducer(config.NsqAddress)
	if err != nil {
		fmt.Printf("init producer failed, err:%v\n", err)
		panic(err)
	}
	Producers[config.TopicName] = producer
	fmt.Println("producer ready ")

	for i := 1; i <= config.ConsumerTotal; i++ {
		i := i
		go func() {
			consumer, err := initConsumer(config.TopicName, config.ChannelName, config.NsqAddress, config.Handler)
			if err != nil {
				fmt.Printf("init consumer failed, err:%v\n", err)
				panic(err)
			}
			fmt.Printf("consumer %v ready", i)
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			<-sigChan
			consumer.Stop()
		}()
	}
}

func GetProducer(topicName string) *nsq.Producer {
	if p, ok := Producers[topicName]; !ok {
		return nil
	} else {
		return p
	}
}

// init producer
func initProducer(str string) (*nsq.Producer, error) {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(str, config)
	if err != nil {
		fmt.Printf("create producer failed, err:%v\n", err)
		return nil, err
	}
	return producer, nil
}

// init consumer
func initConsumer(topicName string, channelName string, address string, handler nsq.Handler) (*nsq.Consumer, error) {
	config := nsq.NewConfig()
	config.LookupdPollInterval = 15 * time.Second
	consumer, err := nsq.NewConsumer(topicName, channelName, config)
	if err != nil {
		fmt.Printf("create consumer failed, err:%v\n", err)
		return nil, err
	}
	consumer.AddHandler(handler)
	//if err = consumer.ConnectToNSQD(address); err != nil { // Direct Connect nsqd
	//	return nil, err
	//}
	if err = consumer.ConnectToNSQLookupd(address); err != nil { // Found via nsqlookupd
		return nil, err
	}
	return consumer, nil
}
