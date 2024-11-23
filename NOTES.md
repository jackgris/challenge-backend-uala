### Additional data:
Database options to PostgreSQL:
https://blog.x.com/engineering/en_us/topics/infrastructure/2023/how-we-scaled-reads-on-the-twitter-users-database

https://discord.com/blog/how-discord-stores-trillions-of-messages

Why xid can be a good ID:
https://encore.dev/blog/go-1.18-generic-identifiers


### I should add this tools:
OAuth: [OAuth2 Proxy](https://github.com/oauth2-proxy/oauth2-proxy)
CI/CD: Github Actions
Gherkin: for E2E test
API documentation: OpenAPI - Swagger
Tracing: OpenTelemetry and Jaeger
Loggin: Elasticsearch and Kibana
Monitoring: Grafana and Prometheus

### Run [migrations](https://github.com/golang-migrate/migrate):

Create tables:
```bash
migrate -database $DATABASE_URL -path ./migrations up
```

Delete tables:
```bash
migrate -database $DATABASE_URL -path ./migrations down
```

Example create a new migration:
```bash
migrate create -ext sql -dir migrations -seq create_likes_table
```

### Run psql (Docker):

```bash
docker exec -it postgresql_tweet_dev psql -d twitter -U pg
```

### Curl test:

#### Health check:
```bash
curl -X GET 'http://localhost:8080/tweet/helthz'
```

#### Create a Tweet example:
```bash
curl -X POST http://localhost:8080/tweet/create \
-H "Content-Type: application/json" \
-d '{"user_id": "1234", "content": "New Sports Event!"}'
```

#### Get Tweet by ID
```bash
curl -X GET 'http://localhost:8080/tweet/id/csuitap82pqc73cn5ar0'
```

#### Delete Tweet
```bash
curl -X DELETE 'http://localhost:8080/tweet/delete/csuitap82pqc73cn5ar0'
```

#### Like
```bash
curl -X POST 'http://localhost:8080/tweet/like' \
-H "Content-Type: application/json" \
-d '{"tweet_id": "csvqhqmek44s73e2qf7g", "user_id": "1234"}'
```

#### Dislike
```bash
curl -X DELETE 'http://localhost:8080/tweet/like' \
-H "Content-Type: application/json" \
-d '{"id":"csvqj46ek44s73e2qf80","tweet_id": "csvqhqmek44s73e2qf7g", "user_id": "1234"}'
```

#### Retweet
```bash
curl -X POST 'http://localhost:8080/tweet/retweet' \
-H "Content-Type: application/json" \
-d '{"tweet_id":"csv5jjrqnq3s73akufl0",  "user_id": "1234"}'
```

#### Remove retweet
```bash
curl -X DELETE 'http://localhost:8080/retweet/csv5jjrqnq3s73akufl0' \
-H "Content-Type: application/json" \
-d '{"id":"csv5ks3qnq3s73akuflg","tweet_id":"csv5jjrqnq3s73akufl0",  "user_id": "1234"}'
```
