
# **Notify Service**

**Notify Service** –  микросервис для обработки сообщений из Kafka, отправки уведомлений. Этот сервис принимает сообщения из топика **`requests`**, обрабатывает их и отправляет результат в топик **`notify-responses`**.

---

## **Содержание**
1. [Описание функциональности](#описание-функциональности)
2. [Архитектура](#архитектура)
3. [Взаимодействие и структура сообщений](#взаимодействие-и-структура-сообщений)
4. [Эндпоинты](#эндпоинты)
5. [Метрики](#метрики)
6. [Запуск и настройка](#запуск-и-настройка)
7. [Конфигурация](#конфигурация)
8. [Пример запуска](#пример-запуска)
9. [Структура проекта](#структура-проекта)
10. [Файлы проекта](#файлы-проекта)

---

## **1. Описание функциональности**

Notify Service выполняет следующие задачи:

- Читает входящие сообщения из Kafka топика **`requests`**.
- Обрабатывает сообщения с определённой структурой.
- Отправляет обработанные ответы в Kafka топик **`notify-responses`**.
- Экспортирует **метрики Prometheus** для мониторинга.
- Логирует все основные этапы работы для дальнейшей отладки и анализа.

---

## **2. Архитектура**

Микросервис состоит из следующих основных компонентов:

- **Kafka Consumer**: Читает входящие сообщения из топика.
- **Kafka Producer**: Публикует ответы в выходной топик.
- **Обработчик сообщений**: Парсит входные сообщения и выполняет бизнес-логику.
- **HTTP-сервер**: Отдаёт метрики Prometheus на эндпоинте `/metrics`.
- **Graceful Shutdown**: Обрабатывает сигналы завершения (`SIGINT`, `SIGTERM`), корректно освобождая ресурсы.

---

## **3. Взаимодействие и структура сообщений**

### **Входные сообщения (Kafka топик: `requests`)**

Формат JSON:

```json
{
  "request_id": "12345",
  "service": "notify",
  "action": "send",
  "data": {
    "user": "example_user",
    "notification": "Welcome to the system!"
  }
}
```

### **Выходные сообщения (Kafka топик: `notify-responses`)**

Формат JSON:

```json
{
  "request_id": "12345",
  "status": "success",
  "data": "Processed action: send"
}
```

---

## **4. Эндпоинты**

| Метод | Путь        | Описание                           |
|-------|-------------|------------------------------------|
| GET   | `/metrics`  | Возвращает метрики Prometheus.     |

---

## **5. Метрики**

Сервис экспортирует метрики в формате **Prometheus** на эндпоинте `/metrics`:

| Метрика                       | Описание                                             |
|-------------------------------|------------------------------------------------------|
| `notify_processed_messages_total` | Общее количество успешно обработанных сообщений.   |
| `notify_error_messages_total`     | Общее количество сообщений, обработка которых завершилась ошибкой. |

### **Пример запроса метрик**

```bash
curl http://localhost:8080/metrics
```

---

## **6. Запуск и настройка**

### **Требования**

- **Go** версии 1.21 или выше.
- **Kafka** (доступный брокер Kafka).
- **Prometheus** (для сбора метрик).

### **1. Установка зависимостей**

```bash
go mod tidy
```

### **2. Конфигурация**

Создайте файл **`.env`** с конфигурацией:

```dotenv
KAFKA_BROKER=localhost:9092
KAFKA_REQUESTS_TOPIC=requests
KAFKA_NOTIFY_RESPONSES_TOPIC=notify-responses
NOTIFY_SERVICE_NAME=notify
LOG_LEVEL=info
```

### **3. Запуск сервиса**

```bash
go run cmd/app/main.go
```

Сервис будет доступен на:

- Метрики: [http://localhost:8080/metrics](http://localhost:8080/metrics)

---

## **7. Конфигурация**

Конфигурация загружается из файла **`.env`** и включает следующие параметры:

| Переменная                     | Описание                                | Пример значения         |
|--------------------------------|-----------------------------------------|-------------------------|
| `KAFKA_BROKER`                 | Адрес Kafka-брокера                    | `localhost:9092`        |
| `KAFKA_REQUESTS_TOPIC`         | Топик для входящих сообщений           | `requests`              |
| `KAFKA_NOTIFY_RESPONSES_TOPIC` | Топик для выходящих сообщений          | `notify-responses`      |

---

## **8. Пример запуска**

1. Убедитесь, что Kafka запущен.
2. Добавьте конфигурацию в `.env`.
3. Запустите сервис:

```bash
go run cmd/notify/main.go
```

4. Проверьте метрики на:

```
http://localhost:8080/metrics
```

---

## **9. Структура проекта**

```
project-root/
├── cmd/
│   └── app/
│       └── main.go            # Основная точка входа
├── config/
│   └── config.go              # Загрузка конфигурации
├── internal/
│   ├── kafka/
│   │   ├── consumer.go        # Kafka Consumer
│   │   ├── producer.go        # Kafka Producer
│   │   └── metrics.go         # Метрики Prometheus
│   ├── handler/
│   │   └── handler.go         # Обработчик сообщений
│   └── logger/
│       └── logger.go          # Логгер
└── go.mod                     # Модуль Go
```

---

## **10. Файлы проекта**

### **main.go** (точка входа)

Содержит логику инициализации Kafka, запуск HTTP-сервера и graceful shutdown.

### **consumer.go**

Обрабатывает чтение сообщений из Kafka.

### **producer.go**

Отправляет сообщения в выходной Kafka-топик.

### **metrics.go**

Экспортирует метрики Prometheus.

### **handler.go**

Обрабатывает входящие сообщения и выполняет бизнес-логику.

### **logger.go**

Инициализирует логгер для структурированного логирования.

### **config.go**

Загружает конфигурацию из файла `.env`.


