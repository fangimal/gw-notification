package storage

import (
	"context"
	"gw-notification/pkg/models"
)

// NotificationStorage определяет контракт для сохранения уведомлений о крупных переводах
type NotificationStorage interface {
	SaveNotification(ctx context.Context, notification *models.Notification) error
}
