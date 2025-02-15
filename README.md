
# **Lost Search**

–  сервис для поиска потерянных вещей
---



## **1. Описание функциональности**

Notify Service выполняет следующие задачи:

- Читает входящие сообщения из Kafka топика **`request`**.
- Обрабатывает сообщения с определённой структурой.
- Отправляет обработанные ответы в Kafka топик **`response`**.
- Экспортирует **метрики Prometheus** для мониторинга.
- Логирует все основные этапы работы для дальнейшей отладки и анализа.

---

## **2. Архитектура**

Микросервис состоит из следующих основных компонентов:

- **Gateway**: Ожидает сообщения на эндпоинтах.
- **Auth-Service**: Авторизация и аутентификация.
- **Ads-Service**: Данные об объявлениях.
- **Notify-Service**: Нотификация пользователей.

---

## **3. Взаимодействие и структура сообщений**

### **Входные сообщения (Kafka топик: `*_requests`)**

Формат JSON:

```json
{
  "request_id": "12345",
  "service": "notify",
  "action": "send",
  "data": {}
}
```

### **Выходные сообщения (Kafka топик: `*_responses`)**

Формат JSON:

```json
{
  "request_id": "12345",
  "status": "success",
  "data": {}
}
```

---

## **4. Эндпоинты**

| Метод | Путь               | Описание             |
|-------|--------------------|----------------------|
| POST  | `/v1/api/Register` | Зарегестрироваться   |
| POST  | `/v1/api/Login`    | Получить новый токен |
| GET   | `/v1/api/GetAds`   | Найти объявление     |
| POST  | `/v1/api/MakeAds`  | Создать объявление   |
| POST  | `/v1/api/ApplyAds` | Откликнуться         |


---


## **5. Структура запросов**


Register:

```json
{
  "login": "",
  "password": "",
  "email": ""
}
```

Login:

```json
{
  "login": "",
  "password": ""
}
```

MakeAds:

```json
{
  "name": "",
  "description": "",
  "type": "",
  "location": {
    "country": "",
    "city": "",
    "district": ""
  }
}
```

SearchAds:

```json
{
  "name": "",        //необязательно
  "description": "", //необязательно
  "type": "",       //необязательно
  "location": {
    "country": "",   //необязательно
    "city": "",    //необязательно
    "district": "" //необязательно
  }
}
```

ApplyAds:

```json
{
  "uuid": ""
}
```


## **6. Структура ответов**


Register:

```json
{
  "success": true/false
}
```

Login:

```json
{
  "success": true/false
}
```

MakeAds:

```json
{
  "name": "",
  "description": "",
  "type": "",
  "location": {
    "country": "",
    "city": "",
    "district": ""
  }
}
```

SearchAds:

```json
{
  "name": "",        //необязательно
  "description": "", //необязательно
  "type": "",       //необязательно
  "location": {
    "country": "",   //необязательно
    "city": "",    //необязательно
    "district": "" //необязательно
  }
}
```

ApplyAds:

```json
{
  "uuid": ""
}
```


