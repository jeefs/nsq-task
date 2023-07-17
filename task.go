package nsqQueue

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

// 注册任务
func sendTask(config TaskConfig) {
	producer, err := initProducer(config.NsqAddress)
	if err != nil {
		fmt.Printf("init producer failed, err:%v\n", err)
		panic(err)
	}
	Producers[config.TopicName] = producer
	fmt.Println("生产者就绪")

	for i := 1; i <= config.ConsumerTotal; i++ {
		i := i
		go func() {
			consumer, err := initConsumer(config.TopicName, config.ChannelName, config.NsqAddress, config.Handler)
			if err != nil {
				fmt.Printf("init consumer failed, err:%v\n", err)
				panic(err)
			}
			fmt.Printf("消费者%v就绪", i)
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

// 初始化生产者
func initProducer(str string) (*nsq.Producer, error) {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(str, config)
	if err != nil {
		fmt.Printf("create producer failed, err:%v\n", err)
		return nil, err
	}
	return producer, nil
}

// 初始化消费者
func initConsumer(topicName string, channelName string, address string, handler nsq.Handler) (*nsq.Consumer, error) {
	config := nsq.NewConfig()
	//config.MaxInFlight = 2
	config.LookupdPollInterval = 15 * time.Second
	consumer, err := nsq.NewConsumer(topicName, channelName, config)
	if err != nil {
		fmt.Printf("create consumer failed, err:%v\n", err)
		return nil, err
	}
	consumer.AddHandler(handler)
	if err = consumer.ConnectToNSQD(address); err != nil { // 直接连NSQD
		return nil, err
	}
	//if err = consumer.ConnectToNSQLookupd(address); err != nil { // 通过lookupd查询
	//	return nil, err
	//}
	return consumer, nil
}
