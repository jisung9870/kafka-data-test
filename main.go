package main

import (
	"flag"
	"random-data/handler"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	brokers   string
	mode      string
	topic     string
	partition int
	count     int
	dataType  string
)

func main() {
	// 연결하기 위한 broker 주소
	flag.StringVar(&brokers, "brokers", "127.0.0.1:9092", "")
	// 어플리케이션 실행 모드를 선택
	flag.StringVar(&mode, "mode", "", "")
	flag.StringVar(&topic, "topic", "test", "")
	flag.IntVar(&partition, "partition", 1, "")
	// producer mode 시 전송하기 위한 메시지 수 - 0인 경우 loop
	// consumer mode 시 연결되는 consumer 수 - 0인 경우 partion 수만큼 생성
	flag.IntVar(&count, "count", 0, "0: loop")
	// message 데이터의 유형을 선택
	flag.StringVar(&dataType, "data-type", "user", "")
	flag.Parse()

	logrus.Infof("Starting mode: %s, broker address: %s\n", mode, brokers)

	brokerList := strings.Split(brokers, ",")

	var wg sync.WaitGroup
	switch mode {
	case "create":
		wg.Add(1)
		go handler.CreateTopic(&wg, brokerList, topic, partition)
	case "producer":
		wg.Add(1)
		go handler.ProducerHandler(&wg, brokerList, topic, count, dataType)
	case "consumer":
		routinCount := count
		if routinCount == 0 {
			// partition 개수와 동일한 consumer gorution 생성
			var err error
			routinCount, err = handler.GetPartition(brokerList, topic)
			if err != nil {
				logrus.Fatalf("error get partition - %v\n", err)
			}
		}
		wg.Add(routinCount)
		for i := 0; i < routinCount; i++ {
			go handler.ConsumerHandler(&wg, brokerList, topic, "consumer-test")
		}
	default:
		logrus.Fatalf("consumer | producer | create - %s\n", mode)
	}
	wg.Wait()
}
