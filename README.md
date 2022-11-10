# Milky Mailer
Simple AMQP based mailer on Golang

## Docker Quick Start
Build the image:
```bash
docker build ./ -t milky-mailer
```
And run it:
```bash
docker run milky-mailer
```

## Configuration

### Environment variables and flags
Milky Mailer uses environment variables and flags to get configuration from `Consul kv`.

Environment variables have higher priority than flags.

| Variable         | Description                                      | Default          |
|------------------|--------------------------------------------------|------------------|
| `APP_NAME`       | `Consul kv` prefix for application configuration | `milky-mailer`   |
| `CONSUL_ADDRESS` | `Consul` address                                 | `localhost:8500` |

### Consul kv
Milky Mailer uses `Consul kv` to store configuration.

All values **required** for work.

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
Milky Mailer support SMTP connection **only over TLS**.

In the future, it will be possible to disable TLS.

## AMQP
Milky Mailer uses `AMQP` to get messages from queue.

### Message headers

| Header        | Description                      | Example                     |
|---------------|----------------------------------|-----------------------------|
| `To`          | E-mail address for `To` header   | `to.user@example.com`       |
| `FromName`    | Name of sender for `From` header | `Milkteam corp.`            |
| `Subject`     | E-mail subject                   | `Hello, world!`             |

### Message body
Message body is a `string` of type `ContentType`.

### Message `ContentType`
Message `ContentType` is a MIMO type (example: `text/plain` or `text/html`). It field uses for email header `ContentType`. 

## TODO list
- [ ] Add support non-TLS SMTP connection
- [ ] Add support for configuration without `Consul kv`
- [ ] Add beautiful error handler
- [ ] Add validation for messages from AMQP
- [ ] Add tests
- [ ] Use `viper` for configuration
- [ ] Replace old `amqp` package to `amqp091-go`

## License
Created by [MilkHunters team](https://milkhunters.ru) under [Creative Commons Zero v1.0 Universal](https://creativecommons.org/publicdomain/zero/1.0/) license.
