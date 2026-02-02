package kafka

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"gw-notification/pkg/logging"
	"gw-notification/pkg/models"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage представляет собой мок-объект для хранилища уведомлений
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) SaveNotification(ctx context.Context, n *models.Notification) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

func TestProcessMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		message        kafka.Message
		expectedError  bool
		expectedUserID int64
	}{
		{
			name: "valid notification",
			message: kafka.Message{
				Value: func() []byte {
					n := &models.Notification{
						UserID:    12345,
						Amount:    100.50,
						Currency:  "USD",
						Timestamp: time.Now(),
					}
					data, _ := json.Marshal(n)
					return data
				}(),
			},
			expectedError:  false,
			expectedUserID: 12345,
		},
		{
			name: "invalid json",
			message: kafka.Message{
				Value: []byte(`invalid json`),
			},
			expectedError: true,
		},
		{
			name: "empty message",
			message: kafka.Message{
				Value: []byte{},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)

			if !tt.expectedError {
				mockStorage.On("SaveNotification", mock.Anything, mock.AnythingOfType("*models.Notification")).Return(nil)
			}

			logger := logging.GetLogger()

			err := ProcessMessage(context.Background(), tt.message, mockStorage, logger)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				mockStorage.AssertExpectations(t)

				// Проверяем, что SaveNotification был вызван с правильным уведомлением
				calls := mockStorage.Calls
				if len(calls) > 0 {
					args := calls[0].Arguments
					if notification, ok := args.Get(1).(*models.Notification); ok {
						assert.Equal(t, tt.expectedUserID, notification.UserID)
					}
				}
			}
		})
	}
}

// Тест для проверки корректного сохранения уведомления
func TestCorrectNotificationSave(t *testing.T) {
	t.Parallel()

	// Подготовка тестовых данных
	testNotification := &models.Notification{
		UserID:    987654321,
		Amount:    250.75,
		Currency:  "EUR",
		Timestamp: time.Date(2023, 10, 15, 14, 30, 0, 0, time.UTC),
	}

	jsonData, err := json.Marshal(testNotification)
	assert.NoError(t, err)

	message := kafka.Message{
		Value: jsonData,
	}

	// Создание мок-объекта хранилища
	mockStorage := new(MockStorage)
	mockStorage.On("SaveNotification", mock.Anything, testNotification).Return(nil)

	logger := logging.GetLogger()

	// Выполнение тестируемой функции
	err = ProcessMessage(context.Background(), message, mockStorage, logger)

	// Проверка результатов
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)

	// Проверка, что метод SaveNotification был вызван с правильными параметрами
	calls := mockStorage.Calls
	assert.Len(t, calls, 1)

	args := calls[0].Arguments
	assert.Equal(t, 2, len(args)) // контекст и уведомление

	ctxArg := args.Get(0).(context.Context)
	assert.NotNil(t, ctxArg)

	notificationArg := args.Get(1).(*models.Notification)
	assert.Equal(t, testNotification.UserID, notificationArg.UserID)
	assert.Equal(t, testNotification.Amount, notificationArg.Amount)
	assert.Equal(t, testNotification.Currency, notificationArg.Currency)
	assert.Equal(t, testNotification.Timestamp, notificationArg.Timestamp)
}
