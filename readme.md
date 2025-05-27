# notifier-app

O notifier app, é uma aplicação resiliente e escalável feita para enviar alertas em diferentes tipos de canais, com facil integração e tolerante a falhas.

# Arquitetura

![Architecture](https://i.imgur.com/MeWpSLM.gif)


# Rotas

## `api/v1/notify` (POST)
> Input
```json
{
  "uuid": "c3a4b7e0-1234-4d5f-abc0-56789abcd123",
  "message": "Você recebeu um pagamento de R$ 90,00!",
  "groups": ["vendas", "marketing"],
  "targetIds": ["slack-hook-123", "discord-hook-999"],
  "event": {
    "name": "payment_success",
    "timestamp": "1748190489",
    "cost_cents": "90000",
    "currency": "BRL",
    "requester": "requester-123@gmail.com",
    "receiver": "guester-123@gmail.com",
    "category": "pix"
  }
}
```

> Output
```json
{
  "message": "Sua mensagem <uuid> foi enviada com sucesso."
}
```
---

## `api/v1/token` (POST)
```json
{
  "admin_user": "James Sergio",
  "key": "environemt-root-key",
}
```

> Output
```json
{
  "message": "Token created"
}
```
---

## `api/v1/:uuid` (GET)
> Output
```json
{
  "uuid": "c3a4b7e0-1234-4d5f-abc0-56789abcd123",
  "channels": ["email", "slack"],
  "event": {
    "name": "payment_success",
    "timestamp": "1748190489",
    "cost_cents": "90000",
    "currency": "BRL",
    "requester": "requester-123@gmail.com",
    "receiver": "guester-123@gmail.com",
    "category": "pix"
  },
  "metadata": {
    "retries": 3,
    "error": "discord failed error, slack failed error, email failed error",
  }
}
```


