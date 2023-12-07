# Mebrox
Mebrox is a message broker server build with only Go's standard library. It uses HTTP and [Server-sent event](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events) for communication.

## Getting started
Clone repository
```bash
git clone git@github.com:meeron/mebrox.git
```

Install dependencies
```bash
go install
```

Test
```bash
make test
```

Build and run
```bash
make run
```

## Api
Default message lock time is 5 minutes. After 3 attempts to consume message, it will be moved to dead letter queue.
### Create topic
```http request
POST http://localhost:3000/topics/<name>
```
### Publish message to topic
```http request
POST http://localhost:3000/topics/<topic_name>/messages
```
### Create subscription
```http request
POST http://localhost:3000/topics/<topic_name>/subscriptions/<subscription_name>
```
### Subscribe
```http request
GET http://localhost:3000/topics/<topic_name>/subscriptions/<subscription_name/subscribe
```
### Commit message
```http request
POST http://localhost:3000/topics/<topic_name>/subscriptions/<subscription_name>/messages/<message_id>/commit
```
