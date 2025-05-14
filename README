# Simple Messaging


## Dependencies

There's only postgresql dependecy, i couldn't find time to implement redis cache. So get it started simply run following command.
```
docker compose up -d
```

If you want to use traces and metrics via open telemetry.simply clone following repository

```
git clone https://github.com/grafana/docker-otel-lgtm
cd docker-otel-lgtm
./run-lgtm.sh
```
Grafana should be available under http://localhost:3000 after a while.

### Migration & Seeding

If you didn't change anything, simply run, it should migrate right away for sample database.
```
make migrate
```

If you want to seed random data to database to be able to process later.

```
make seed
```

### Webhook Site

If you want to change webhook site to custom one (or current one expired)

change following config under .env file

```
WEBHOOK_SITE_URL=https://webhook.site/8f55092f-5110-4f06-896c-0527fb97e8b0
```

dont forget to change reponse code to 202 and data to following

```
{"message": "Accepted","messageId":"03ca588e-8269-4871-8990-4ecc23e8e967"}
```

### Updating Docs

Swagger files generated under docs folder, if you want to change something and update docs
```
make docs
```

## Starting Worker

Simply run the following command.

```
make run
```

server should be available under http://localhost:8080

Sample urls are (details are in the swagger):

GET http://localhost:8080/api/v1/messages/list?page=1

POST http://localhost:8080/api/v1/service/toggle
