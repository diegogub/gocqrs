package api

import (
	"bitbucket.org/dgub/evento/dom"
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"log"
)

var (
	nsqConsumer *nsq.Consumer
)

func (cli *ReaderClient) connectNSQ(nsqServer string) error {
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("evento", cli.Id, config)
	if err != nil {
		return err
	}

	nsqConsumer = consumer
	nsqConsumer.AddHandler(cli)
	err = nsqConsumer.ConnectToNSQD(nsqServer)
	if err != nil {
		return err
	}
	return nil
}

func (rcli *ReaderClient) HandleMessage(msg *nsq.Message) error {
	var event dom.Event
	err := json.Unmarshal(msg.Body, &event)

	if DevMode {
		log.Println("Event: ", event)
		log.Println("Error : ", err)
	}

	if err != nil {
		return err
	}

	err = event.Validate()
	if err != nil {
		return err
	}

	if rcli.StreamMatch(event.StreamId) {
		// Send Event
		rcli.wg.Add(1)
		rcli.chEvents <- &event
	}
	msg.Finish()
	return nil
}
