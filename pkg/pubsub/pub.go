package pubsub

import (
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/google/uuid"
)

type Publisher struct {
	messages  []*message.Message
	publisher *amqp.Publisher
}

func NewPublisher() (*Publisher, error) {
	amqpConfig, err := DurablePubSubConfig()
	if err != nil {
		return nil, err
	}
	publisher, err := amqp.NewPublisher(
		*amqpConfig,
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		messages:  make([]*message.Message, 0),
		publisher: publisher,
	}, nil
}

func (pub *Publisher) Update(mid string, uid, rid, unid *uuid.UUID, body interface{}) error {
	byteMsg, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_uid := uuid.New()
	if uid != nil {
		_uid = *uid
	}

	sendMsg := Msg{
		MsgBase: MsgBase{
			MID:    mid,
			UID:    _uid,
			Sender: Sender(),
			RID:    rid,
			UnID:   unid,
		},
		Body: string(byteMsg),
	}

	sendByteMsg, err := json.Marshal(sendMsg)
	if err != nil {
		return err
	}

	pub.messages = append(
		pub.messages,
		message.NewMessage(
			_uid.String(),
			sendByteMsg,
		),
	)

	return nil
}

func (pub *Publisher) Publish() error {
	return pub.publisher.Publish(
		GlobalPubsubTopic,
		pub.messages...,
	)
}

func (pub *Publisher) Close() {
	pub.publisher.Close()
}

func WithPublisher(updater func(publisher *Publisher) error) error {
	if updater == nil {
		return fmt.Errorf("invalid updater")
	}

	publisher, err := NewPublisher()
	if err != nil {
		return err
	}
	defer publisher.Close()

	if err := updater(publisher); err != nil {
		return err
	}

	return publisher.Publish()
}
