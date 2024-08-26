package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"kafka-sarama-example/logger"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"go.uber.org/zap"
	"gocloud.dev/pubsub"
	kpb "gocloud.dev/pubsub/kafkapubsub"
)

type Subscriber struct {
	sub    *pubsub.Subscription
	logger *zap.Logger
}

func NewSubscriber(brokers []string, topic, groupID string) (*Subscriber, error) {
	config := kpb.MinimalConfig()
	sub, err := kpb.OpenSubscription(
		brokers,
		config,
		groupID,
		[]string{topic},
		&kpb.SubscriptionOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("kafkapubsub.OpenSubscription: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %v", err)
	}

	return &Subscriber{
		sub:    sub,
		logger: logger,
	}, nil
}

func (s *Subscriber) Start(ctx context.Context) error {
	for {
		msg, err := s.sub.Receive(ctx)
		if err != nil {
			return fmt.Errorf("error receiving message: %v", err)
		}
		go s.handleMessage(msg)
	}
}

func (s *Subscriber) handleMessage(msg *pubsub.Message) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("パニック発生", zap.Any("error", r), zap.String("message", string(msg.Body)))
		}
		s.logger.Info("Ack開始")
		msg.Ack()
		s.logger.Info("Ack完了")
	}()

	s.logger.Info("Received message",
		zap.String("data", string(msg.Body)),
		zap.Any("metadata", msg.Metadata),
	)

	var data map[string]interface{}
	if err := json.Unmarshal(msg.Body, &data); err != nil {
		s.logger.Error("Failed to decode JSON", zap.Error(err))
		panic(err)
		return
	}

	if err := s.processMessage(data); err != nil {
		s.logger.Error("Failed to process message", zap.Error(err))
		return
	}

	s.logger.Info("Processed message", zap.Any("data", data))
}

func (s *Subscriber) processMessage(data map[string]interface{}) error {
	// Process the message here
	// For example: save to database, notify other services, etc.
	s.logger.Info("Processing message", zap.Any("data", data))
	return nil
}

func (s *Subscriber) Close() error {
	return s.sub.Shutdown(context.Background())
}

func Run(brokers, version, topic string) {
	if brokers == "" {
		logger.Fatal("At least one broker is required")
	}
	splitBrokers := strings.Split(brokers, ",")
	subscriber, err := NewSubscriber(splitBrokers, topic, "default_group")
	if err != nil {
		logger.Fatal("failed to create subscriber", zap.Error(err))
	}
	defer subscriber.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// シグナル処理を設定
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		errCh <- subscriber.Start(ctx)
	}()

	select {
	case err := <-errCh:
		logger.Error("subscriber error", zap.Error(err))
	case sig := <-sigCh:
		logger.Info("Received signal, shutting down", zap.String("signal", sig.String()))
		cancel()
	}


	logger.Sync()
}