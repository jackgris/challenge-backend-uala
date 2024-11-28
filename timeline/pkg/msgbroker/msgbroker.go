package msgbroker

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackgris/twitter-backend/timeline/pkg/logger"
	"github.com/jackgris/twitter-backend/timeline/pkg/uuid"
	nc "github.com/nats-io/nats.go"
)

type MsgBroker struct {
	sub  message.Subscriber
	pub  message.Publisher
	logs *logger.Logger
}

func NewMsgBroker(path string, logs *logger.Logger) *MsgBroker {

	marshaler := &nats.GobMarshaler{}
	logger := watermill.NewStdLogger(false, false)
	options := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30 * time.Second),
		nc.ReconnectWait(1 * time.Second),
	}
	jsConfig := nats.JetStreamConfig{Disabled: true}

	conn, err := nc.Connect(path)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
		defer cancel()

		logs.Error(ctx, "timeline service", "Failed to connect to Nats server", err)
	}

	js, err := conn.JetStream()
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
		defer cancel()

		logs.Error(ctx, "timeline service", "Failed to get JetStream context", err)
	}

	streamConfig := &nc.StreamConfig{
		Name:      "tweet-stream",
		Subjects:  []string{"tweet.*"},
		Retention: nc.LimitsPolicy,
		Storage:   nc.FileStorage, // or nats.MemoryStorage
		Replicas:  1,              // Number of replicas
	}

	_, err = js.AddStream(streamConfig)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
		defer cancel()

		logs.Error(ctx, "timeline service", "Failed to get add Stream", err)
	}

	configPub := nats.PublisherConfig{
		URL:         path,
		NatsOptions: options,
		Marshaler:   marshaler,
		JetStream:   jsConfig,
	}

	publisher, err := nats.NewPublisherWithNatsConn(conn, configPub.GetPublisherPublishConfig(), logger)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
		defer cancel()

		logs.Error(ctx, "timeline service", "publisher: connection to message broker", "status", err)
	}

	configSub := nats.SubscriberConfig{
		URL:         path,
		NatsOptions: options,
		JetStream:   jsConfig,
	}

	subscriber, err := nats.NewSubscriberWithNatsConn(conn, configSub.GetSubscriberSubscriptionConfig(), logger)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
		defer cancel()

		logs.Error(ctx, "timeline service", "subscriber: connection to message broker", "status", err)
	}

	return &MsgBroker{pub: publisher, sub: subscriber, logs: logs}
}

func (m *MsgBroker) PublishMessages(topic string, msg *message.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()
loop:
	for {
		select {
		case <-ctx.Done():
			m.logs.Info(ctx, "timeline service", "publish message", "status", "conxtext timeout", "message ID", msg.UUID)
			break loop
		default:
			if err := m.pub.Publish(topic, msg); err != nil {
				m.logs.Info(ctx, "timeline service", "publish message", "status", err, "message ID", msg.UUID)
			} else {
				break loop
			}
		}
		time.Sleep(time.Millisecond * 50)
	}
}

func (m *MsgBroker) SubscribeGetFollowers(ctx context.Context) {

	messages, err := m.sub.Subscribe(ctx, "get_followers")
	if err != nil {
		m.logs.Error(ctx, "timeline service", "subscriber: can't subscribe get_followers", "status", err)
	}

	go func() {
		for msg := range messages {
			m.logs.Info(ctx, "timeline service", "Message ID get followers", msg.UUID, string(msg.Payload))
		}
	}()
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
