package kafka

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"gw-notification/pkg/logging"
	"gw-notification/pkg/models"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/mock"
)

// MockStorageBench представляет собой мок-объект для хранилища уведомлений для бенчмарков
type MockStorageBench struct {
	mock.Mock
}

func (m *MockStorageBench) SaveNotification(ctx context.Context, n *models.Notification) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

// Бенчмарк для проверки производительности - обработка до 1000 сообщений в секунду
func BenchmarkNotificationProcessing(b *testing.B) {
	b.ReportAllocs()

	// Подготовка тестовых данных
	testNotification := &models.Notification{
		UserID:    123456789,
		Amount:    150.25,
		Currency:  "USD",
		Timestamp: time.Now(),
	}

	jsonData, _ := json.Marshal(testNotification)

	message := kafka.Message{
		Value: jsonData,
	}

	// Создание мок-объекта хранилища
	mockStorage := new(MockStorageBench)
	mockStorage.On("SaveNotification", mock.Anything, mock.AnythingOfType("*models.Notification")).Return(nil)

	// Создаем логгер для бенчмарка
	loggingBenchLogger := logging.GetLogger()

	b.ResetTimer()

	// Запуск бенчмарка
	for i := 0; i < b.N; i++ {
		_ = ProcessMessage(context.Background(), message, mockStorage, loggingBenchLogger)
	}

	// Подсчет количества операций в секунду
	opsPerSec := float64(b.N) / b.Elapsed().Seconds()
	b.ReportMetric(opsPerSec, "ops/sec")

	// Проверка, достигается ли целевая производительность
	targetOpsPerSec := 1000.0
	if opsPerSec >= targetOpsPerSec {
		loggingBenchLogger.Infof("✓ Производительность: %.2f операций/сек (цель: %d операций/сек)", opsPerSec, int(targetOpsPerSec))
	} else {
		loggingBenchLogger.Infof("✗ Производительность: %.2f операций/сек (цель: %d операций/сек)", opsPerSec, int(targetOpsPerSec))
	}
}

// Дополнительный бенчмарк для тестирования обработки большого количества сообщений
func BenchmarkMassiveNotificationProcessing(b *testing.B) {
	b.ReportAllocs()

	// Подготовка различных тестовых данных
	notifications := make([]*models.Notification, 100)
	for i := 0; i < 100; i++ {
		notifications[i] = &models.Notification{
			UserID:    int64(100000000 + i),
			Amount:    float32(100+i) + 0.99,
			Currency:  "USD",
			Timestamp: time.Now().Add(time.Duration(i) * time.Second),
		}
	}

	messages := make([]kafka.Message, 100)
	for i, notification := range notifications {
		jsonData, _ := json.Marshal(notification)
		messages[i] = kafka.Message{
			Value: jsonData,
		}
	}

	// Создание мок-объекта хранилища
	mockStorage := new(MockStorageBench)
	mockStorage.On("SaveNotification", mock.Anything, mock.AnythingOfType("*models.Notification")).Return(nil).Maybe()

	// Создаем логгер для бенчмарка
	loggingBenchLogger := logging.GetLogger()

	b.ResetTimer()

	// Запуск бенчмарка
	for i := 0; i < b.N; i++ {
		msg := messages[i%len(messages)] // циклически используем сообщения
		_ = ProcessMessage(context.Background(), msg, mockStorage, loggingBenchLogger)
	}

	// Подсчет количества операций в секунду
	opsPerSec := float64(b.N) / b.Elapsed().Seconds()
	b.ReportMetric(opsPerSec, "ops/sec")

	// Проверка, достигается ли целевая производительность
	targetOpsPerSec := 1000.0
	if opsPerSec >= targetOpsPerSec {
		loggingBenchLogger.Infof("✓ Массовая обработка: %.2f операций/сек (цель: %d операций/сек)", opsPerSec, int(targetOpsPerSec))
	} else {
		loggingBenchLogger.Infof("✗ Массовая обработка: %.2f операций/сек (цель: %d операций/сек)", opsPerSec, int(targetOpsPerSec))
	}
}
