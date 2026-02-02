package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"gw-notification/iternal/config"
	"gw-notification/iternal/storage"
	mongodb "gw-notification/iternal/storage/mongo"
	"gw-notification/pkg/logging"
	"gw-notification/pkg/models"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(cfg *config.Config, logger *logging.Logger, db *mongodb.Storage) error {
	// init Kafka conf
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{cfg.Kafka.KafkaServerAddress},
		Topic:     cfg.Kafka.KafkaTopic,
		GroupID:   cfg.Kafka.KafkaGroupID,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		Partition: 0,
	})
	defer kafkaReader.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to signal when to stop
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := kafkaReader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				logger.Error("Failed to read message", "error", err)
				continue
			}

			var notification models.Notification
			if err = json.Unmarshal(msg.Value, &notification); err != nil {
				logger.Error("Failed to unmarshal message", "error", err)
				continue
			}

			if err = db.SaveNotification(ctx, &notification); err != nil {
				logger.Error("Failed to save notification", "error", err)
				continue
			}

			logger.WithFields(map[string]interface{}{
				"user_id":  notification.UserID,
				"amount":   notification.Amount,
				"currency": notification.Currency,
			}).Info("Notification saved")

			if err = kafkaReader.CommitMessages(ctx, msg); err != nil {
				logger.Warn("Failed to commit offset", "error", err)
			}
		}
	}
}

func ProcessMessage(ctx context.Context, msg kafka.Message, storage storage.NotificationStorage, logger *logging.Logger) error {
	var notification models.Notification
	if err := json.Unmarshal(msg.Value, &notification); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	if err := storage.SaveNotification(ctx, &notification); err != nil {
		return fmt.Errorf("save: %w", err)
	}

	//logger.Info("Notification saved", "user_id", notification.UserID)
	return nil
}
