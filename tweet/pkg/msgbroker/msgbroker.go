package msgbroker

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackgris/twitter-backend/tweet/pkg/logger"
	"github.com/jackgris/twitter-backend/tweet/pkg/uuid"
	nc "github.com/nats-io/nats.go"
)

type MsgBroker struct {
	name string
	sub  message.Subscriber
	pub  message.Publisher
	logs *logger.Logger
}

func NewMsgBroker(service, natsURL string, logs *logger.Logger) *MsgBroker {

	marshaler := &nats.GobMarshaler{}
	logger := watermill.NewStdLogger(false, false)
	options := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30 * time.Second),
		nc.ReconnectWait(1 * time.Second),
	}
	subscribeOptions := []nc.SubOpt{
		nc.DeliverAll(),
		nc.AckExplicit(),
	}

	jsConfig := nats.JetStreamConfig{
		Disabled:         false,
		AutoProvision:    true,
		ConnectOptions:   nil,
		SubscribeOptions: subscribeOptions,
		PublishOptions:   nil,
		TrackMsgId:       false,
		AckAsync:         false,
		DurablePrefix:    "",
	}

	conn, err := nc.Connect(natsURL)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
		defer cancel()

		logs.Error(ctx, service, "Failed to connect to Nats server", err)
	}

	configPub := nats.PublisherConfig{
		URL:         natsURL,
		NatsOptions: options,
		Marshaler:   marshaler,
		JetStream:   jsConfig,
	}

	publisher, err := nats.NewPublisher(configPub, logger)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
		defer cancel()

		logs.Error(ctx, service, "publisher: connection to message broker", "status", err)
	}

	configSub := nats.SubscriberConfig{
		URL:            natsURL,
		CloseTimeout:   30 * time.Second,
		AckWaitTimeout: 30 * time.Second,
		NatsOptions:    options,
		Unmarshaler:    marshaler,
		JetStream:      jsConfig,
	}

	subscriber, err := nats.NewSubscriberWithNatsConn(conn, configSub.GetSubscriberSubscriptionConfig(), logger)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
		defer cancel()

		logs.Error(ctx, service, "subscriber: connection to message broker", "status", err)
	}

	return &MsgBroker{name: service, pub: publisher, sub: subscriber, logs: logs}
}

func (m *MsgBroker) PublishMessages(topic string, msg *message.Message) {
	if err := m.pub.Publish(topic, msg); err != nil {
		m.logs.Info(context.Background(), m.name, "publish message", "status", err, "topic", topic, "message ID", msg.UUID)
	}
}

func (m *MsgBroker) SubscribeEvents(topic string) (<-chan *message.Message, error) {
	ctx := context.Background()
	messages, err := m.sub.Subscribe(ctx, topic)
	if err != nil {
		m.logs.Error(ctx, m.name, "subscriber: can't subscribe to "+topic, "status", err)
		return nil, err
	}

	return messages, err
}

func (m *MsgBroker) Close() {
	m.sub.Close()
}

type Header struct {
	ID            string `json:"id"`
	EventName     string `json:"event_name"`
	CorrelationID string `json:"correlation_id"`
	PublishedAt   string `json:"published_at"`
}

func NewHeader(eventName string) Header {
	return Header{
		ID:          uuid.New(),
		EventName:   eventName,
		PublishedAt: time.Now().Format(time.RFC3339),
	}
}

type MockPublisher struct{}

func (m *MockPublisher) Publish(topic string, messages ...*message.Message) error {
	return nil
}

func (m *MockPublisher) Close() error {
	return nil
}

type MockSubscriber struct{}

func (m *MockSubscriber) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	return nil, nil
}

func (m *MockSubscriber) Close() error {
	return nil
}

func NewMockMsgBroker(logs *logger.Logger) *MsgBroker {
	return &MsgBroker{
		pub:  &MockPublisher{},
		sub:  &MockSubscriber{},
		logs: logs,
	}
}
