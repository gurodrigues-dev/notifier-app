# Notifier App

This document outlines the architecture and usage of the Notifier App, which consists of two services: an HTTP API and a Kafka consumer. The architecture is detailed below.

To send a payment notification via the API, a token is required for authentication, implemented as a simple mechanism. To generate a new token, a system administrator must use an admin token, which is sent via email. This admin token must be included in the header for endpoints under `/token`.

With a generated token, you can register a channel to receive notifications. When creating a channel, you must specify a group it belongs to, e.g., "development." When sending a notification, you can specify either the channel ID or the group. If a group is specified and multiple channels are registered under it, all channels in the group will receive the notification. For example, sending to `["development", "2", "marketing"]` will notify all channels in the "development" and "marketing" groups, plus the specific channel with ID "2" (which could belong to an admin or another entity). By default, three platforms are supported: Email, Slack, and Discord. The system is designed to decouple the addition of new channel types, making it easy to extend.

When registering a channel:
- For email channels, a confirmation email is sent.
- For Slack or Discord, a `webhookURL` is required (e.g., a Slack app must provide a `webhook_url`).

Each notification includes a UUID in the HTTP request (assumed to align with microservices communication, where service X, the client, calls this API, service Y). The UUID is checked against a cache to detect if the message is already being processed. If it is, the request is acknowledged and returned. The cache prevents message duplication and overload. The message is then sent to a Kafka queue. On success, it’s recorded in the cache; on failure, it’s logged in a database and metrics are incremented in Prometheus, available at the `/metrics` endpoint for instrumentation.

The Kafka consumer processes messages from the queue. Each message includes metadata with a retry count for error handling. The consumer:
1. Checks the retry count. If it exceeds 3, an error is returned.
2. If retries are 1 or 2, the message is requeued in Kafka for reprocessing.
3. Attempts to send the message to the designated platform(s). If multiple platforms are involved, success on at least one platform meets the SLA.
4. If all attempts fail and incrementing retries reaches 3, the message is logged in the error database as an alert, and a Prometheus metric is incremented for monitoring (configurable with Grafana and AlertManager).

For application management, endpoints are provided to:
- View errors in the error database (admin-only).
- Delete tokens (admin token required) or channels (user token required).
- List channels by ID (`channel/9`), group (`group/marketing`), or platform (`platform/discord`).

## 🔬 Developer Notes

### Scalability

To minimize infrastructure dependencies, Kubernetes was not used. Instead, Docker Swarm is configured with replicas for horizontal scaling, providing built-in load balancing. For VM deployments, K3s would typically be used, but Docker Swarm suffices for this challenge.

### Frontend Tools

Frontend tools are provided for local visualization of application data:
- **Metabase**: Accessible at port 3001. Create a fast account and connect to the local database (configured in the Docker network for ease of use). View data from the Channels, Tokens, and Errors tables.
- **Kafka UI**: Accessible at port 8092. Create a fast account and add broker details from the `create-topics` service in the Docker Compose or Swarm configuration.

## 🛠️ Tools Used

- Golang (Gin)
- Kafka
- ZooKeeper
- Redis
- Postgres
- Grafana
- Prometheus
- AlertManager
- Kafka-UI
- Metabase

*Note*: Instrumentation, monitoring, and observability tools require configuration based on the deployment environment.

## 🚀 How to Use

Run the application locally with Docker Compose or in a VM with Docker Swarm.

For local setup:
```bash
docker compose up --build -d
```

The application should start successfully. If database issues arise, a Makefile is included in the repository (copied into the container). To run migrations, access the container or run locally:
```bash
make migrateup
```

This updates the database, making it ready for use.

## ⚙️ Endpoints

### POST /api/v1/token

Create a new user token.

**Parameters**

| Name        | Location | Type   | Description         |
|-------------|----------|--------|---------------------|
| `admin_user`| Body     | String | User's email        |

**Response**

```json
"b23a4b2d-f40c-4627-8d85-89597e67f41d"
```

---

### GET /api/v1/token/:email

Retrieve a user's token by email.

**Parameters**

| Name   | Location | Type   | Description         |
|--------|----------|--------|---------------------|
| `email`| Request  | String | User's email        |

**Response**

```json
{
    "id": 7,
    "token": "5e7f966f-42c0-475d-aa7c-235806cf5391",
    "admin_user": "guester1234@gmail.com"
}
```

---

### DELETE /api/v1/token/:email

Delete a user's token by email.

**Parameters**

| Name   | Location | Type   | Description         |
|--------|----------|--------|---------------------|
| `email`| Request  | String | User's email        |

**Response**

```json
{
    "message": "token deleted successfully"
}
```

---

### POST /api/v1/channel

Create a new channel for notifications.

**Parameters**

| Name        | Location | Type   | Description                     |
|-------------|----------|--------|---------------------------------|
| `platform`  | Body     | String | Platform used (e.g., email, slack, discord) |
| `target_id` | Body     | String | Email or Webhook URL            |
| `group`     | Body     | String | Group name                      |

**Response**

```json
{
    "id": 9,
    "platform": "discord",
    "target_id": "https://discord.com/api/webhooks/1377021857tLnJOk_z",
    "group": "marketing"
}
```

---

### GET /api/v1/channel/:id

Retrieve channel details by ID.

**Parameters**

| Name | Location | Type   | Description         |
|------|----------|--------|---------------------|
| `id` | Request  | String | Channel ID          |

**Response**

```json
{
    "id": 9,
    "platform": "discord",
    "target_id": "https://discord.com/api/webhooks/1377021857500369049/s7AoSlgLs4LCa8qiczRhW9RrODkATsPXJW6REPHngTNM-rBOjuqWISfjmi2qtLnJOk_s",
    "group": "challenge"
}
```

---

### GET /api/v1/group/:group

Retrieve channels by group name.

**Parameters**

| Name    | Location | Type   | Description         |
|---------|----------|--------|---------------------|
| `group` | Request  | String | Group name          |

**Response**

```json
[
    {
        "id": 5,
        "platform": "email",
        "target_id": "joaoteste@gmail.com",
        "group": "development"
    },
    {
        "id": 8,
        "platform": "email",
        "target_id": "gustavo@gmail.com",
        "group": "development"
    }
]
```

---

### GET /api/v1/platform/:platform

Retrieve channels by platform.

**Parameters**

| Name       | Location | Type   | Description         |
|------------|----------|--------|---------------------|
| `platform` | Request  | String | Platform name       |

**Response**

```json
[
    {
        "id": 4,
        "platform": "slack",
        "target_id": "https://hooks.slack.com/services/token",
        "group": "sumup"
    }
]
```

---

### DELETE /api/v1/channel/:id

Delete a channel by ID.

**Parameters**

| Name | Location | Type   | Description         |
|------|----------|--------|---------------------|
| `id` | Request  | String | Channel ID          |

**Response**

```json
{
    "message": "channel deleted successfully"
}
```

---

### POST /api/v1/notification

Create a notification for the specified channels.

**Parameters**

| Name        | Location | Type         | Description                     |
|-------------|----------|--------------|---------------------------------|
| `uuid`      | Body     | String       | Message identifier              |
| `title`     | Body     | String       | Message title                   |
| `channels`  | Body     | Array[String]| Target channels or groups       |
| `event`     | Body     | Map          | Event details                   |
| `name`      | Event    | String       | Notification name               |
| `timestamp` | Event    | Int64        | Notification timestamp          |
| `cost_cents`| Event    | Int64        | Transferred amount (in cents)   |
| `currency`  | Event    | String       | Currency type                   |
| `requester` | Event    | String       | Sender of the transfer          |
| `receiver`  | Event    | String       | Recipient of the transfer       |
| `category`  | Event    | String       | Transfer type                   |

**Example Request**

```json
{
    "uuid": "2bbcdd20-1ea6-42be-8484-02f3007e3473",
    "title": "You received a TED transaction",
    "channels": ["9"],
    "event": {
        "name": "payment_success",
        "timestamp": 1748190489,
        "cost_cents": 90000,
        "currency": "BRL",
        "requester": "requester123@gmail.com",
        "receiver": "guester123@gmail.com",
        "category": "pix"
    }
}
```

**Response**

```json
{
    "message": "notification sent successfully"
}
```

---

### GET /api/v1/notification/:uuid

Retrieve details of a failed notification.

**Parameters**

| Name   | Location | Type   | Description         |
|--------|----------|--------|---------------------|
| `uuid` | Request  | String | Message identifier  |

**Response**

```json
{
    "id": 1,
    "uuid": "2bbcdd20-1ea6-42be-8484-02f3007e3463",
    "body": {
        "id": 0,
        "uuid": "2bbcdd20-1ea6-42be-8484-02f3007e3463",
        "title": "You received a Pix",
        "message": "There was a new transaction at 25/05/25 13:28, between requester123@gmail.com and guester123@gmail.com by pix, with the value of BRL 900, status: payment_success",
        "channels": {
            "1": {
                "id": 1,
                "platform": "discord",
                "target_id": "https://discord.com/api/webhooks/url",
                "group": "marketing"
            },
            "4": {
                "id": 4,
                "platform": "slack",
                "target_id": "https://hooks.slack.com/services/webhook",
                "group": "marketing"
            }
        },
        "event": {
            "name": "payment_success",
            "currency": "BRL",
            "requester": "requester123@gmail.com",
            "receiver": "guester123@gmail.com",
            "category": "pix",
            "timestamp": 1748190489,
            "cost_cents": 90000
        },
        "retries": 3
    },
    "error": "message could not be queued"
}
```

## 📷 Evidence (Slack, Discord, Email)

| Evidence       | Description                          | Preview |
|----------------|--------------------------------------|---------|
| Slack Message  | Notification sent to Slack channel   | ![Slack Image](https://i.imgur.com/0OyuTdf.png) |
| Discord Message| Notification sent to Discord server  | ![Discord Image](https://i.imgur.com/E6pvaa3.png) |
| Email Message  | Notification received via Email      | ![Email Image](https://i.imgur.com/Zcpgw0O.png) |

## 📊 Flowcharts

| Evidence      | Description                       | Preview |
|---------------|-----------------------------------|---------|
| Usage Model   | General usage model of the system | ![Usage Model](https://i.imgur.com/3cBtfo5.png) |
| App Flow      | Internal flow of the application  | ![App Flow](https://i.imgur.com/I8UuSc9.png) |
| Swagger       | Documentation of the application  | ![App Flow](https://i.imgur.com/v8Kv6dx.png) |

## 💡 Architecture

https://github.com/user-attachments/assets/54e660d2-4430-49f9-83db-b84f6b439087
