# Service Notifications

Сервис уведомлений, который получает сообщения о событиях (например, финансовых транзакциях) через Kafka и сохраняет их в MongoDB для дальнейшей обработки и хранения истории операций.

## Описание

Сервис уведомлений (Notification Service) предназначен для приема, обработки и хранения уведомлений о различных событиях в системе. Он использует Apache Kafka в качестве брокера сообщений для получения уведомлений от других сервисов и MongoDB для долгосрочного хранения этих уведомлений.

### Функциональность

1. Получение сообщений из топика Kafka
2. Сохранение уведомлений в MongoDB
3. Обеспечение надежной обработки сообщений
4. Поддержка высокой производительности (до 1000 сообщений в секунду)
5. Структурированное логирование в JSON формате
6. Поддержка graceful shutdown

### Технологии

- **Go** - язык программирования
- **Apache Kafka** - брокер сообщений
- **MongoDB** - NoSQL база данных
- **Docker** - контейнеризация
- **Docker Compose** - оркестрация контейнеров

### Архитектура

Сервис состоит из следующих компонентов:

- **Kafka Consumer** - читает сообщения из топика Kafka
- **MongoDB Storage** - интерфейс для сохранения уведомлений в базе данных
- **Logger** - обеспечивает структурированное логирование
- **Configuration Manager** - загружает и управляет настройками сервиса

### Структура проекта

```
gw-notification/
├── cmd/
│   └── main.go                 # Точка входа в приложение
├── internal/
│   ├── broker/
│   │   └── kafka/
│   │       └── consumer.go     # Обработчик Kafka сообщений
│   ├── config/                 # Конфигурация приложения
│   │   └── config.go
│   ├── storage/
│   │   ├── storage.go          # Интерфейс хранилища
│   │   └── mongo/
│   │       └── mongo.go        # Реализация MongoDB хранилища
├── pkg/
│   ├── models/
│   │   └── models.go           # Структуры данных
│   └── logging/
│       └── logging.go          # Конфигурация логирования
├── logs/
├── config.yml                  # Конфигурационный файл
├── docker-compose.yml          # Docker конфигурация
├── Dockerfile                  # Docker образ
└── go.mod                      # Зависимости Go модуля
```

### Модель уведомления

```go
type Notification struct {
    UserID    int64     `json:"user_id" bson:"user_id"`
    Amount    float32   `json:"amount" bson:"amount"`
    Currency  string    `json:"currency" bson:"currency"`
    Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}
```

### Конфигурация

Сервис использует следующие параметры конфигурации:

- `storage` - строка подключения к MongoDB
- `kafka.ConsumerGroup` - группа потребителей Kafka
- `kafka.ConsumerTopic` - топик Kafka для потребления
- `kafka.ConsumerPort` - порт для соединения с Kafka
- `kafka.KafkaServerAddress` - адрес сервера Kafka
- `kafka.KafkaTopic` - топик Kafka
- `kafka.KafkaGroupID` - ID группы Kafka

### Логирование

Сервис реализует структурированное логирование с поддержкой следующих уровней:
- DEBUG
- INFO
- WARN
- ERROR

Каждое лог-сообщение содержит контекст выполнения (например, user_id, amount, currency).

### Запуск сервиса

#### Через Docker Compose

```bash
docker-compose up --build
```

#### Локальный запуск

1. Убедитесь, что запущены Kafka и MongoDB
2. Установите зависимости: `go mod tidy`
3. Запустите приложение: `go run cmd/main.go`

### Тестирование

Сервис включает unit-тесты и benchmark тесты для проверки производительности:

- `consumer_test.go` - unit-тесты для обработчика Kafka
- `consumer_bench_test.go` - benchmark тесты производительности

### Мониторинг и метрики

Сервис поддерживает мониторинг через Prometheus метрики.

### Безопасность и надежность

- Поддержка graceful shutdown для корректного завершения работы
- Проверка ошибок при обработке сообщений
- Коммит оффсетов Kafka после успешной обработки сообщений
