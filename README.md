# Milky Mailer
Simple AMQP based mailer on Golang

## Docker Quick Start
Build the image:
```
```bash
docker build ./ -t milky-mailer
```
And run it:
```bash
docker run milky-mailer
```

## Configuration

### Environment variables
Milky Mailer uses environment variables to get configuration from `Consul kv`.

| Variable         | Description                                      | Default          |
|------------------|--------------------------------------------------|------------------|
| `APP_NAME`       | `Consul kv` prefix for application configuration | `milky-mailer`   |
| `CONSUL_ADDRESS` | `Consul` address                                 | `localhost:8500` |

### Consul kv
Milky Mailer uses `Consul kv` to store configuration.

All values **required** for the start.

| Key                       | Description                      | Example            |
|---------------------------|----------------------------------|--------------------|
| `APP_NAME/amqp/host`      | AMQP Host                        | `127.0.0.1`        |
| `APP_NAME/amqp/port`      | AMQP Port                        | `5672`             |
| `APP_NAME/amqp/user`      | AMQP User                        | `user`             |
| `APP_NAME/amqp/password`  | AMQP Password                    | `somepassword`     |
| `APP_NAME/amqp/queue`     | AMQP Queue                       | `email`            |
| `APP_NAME/email/host`     | SMTP Host                        | `smtp.example.com` |
| `APP_NAME/email/port`     | SMTP Port                        | `465`              |
| `APP_NAME/email/user`     | SMTP User                        | `user@example.com` |
| `APP_NAME/email/password` | SMTP Password                    | `somepassword`     |
| `APP_NAME/email/from`     | E-mail address for `From` header | `user@example.com` |


## TLS
Milky Mailer support **only TLS** SMTP connection.

## AMQP
Milky Mailer uses `AMQP` to get messages from queue.

### Message headers

| Header        | Description                      | Example                     |
|---------------|----------------------------------|-----------------------------|
| `To`          | E-mail address for `To` header   | `to.user@example.com`       |
| `FromName`    | Name of sender for `From` header | `Milkteam corp.`            |
| `Subject`     | E-mail subject                   | `Hello, world!`             |
| `ContentType` | E-mail content type of body      | `text/plain` or `text/html` |

### Message body
Message body is a `string` of type `ContentType` from headers.