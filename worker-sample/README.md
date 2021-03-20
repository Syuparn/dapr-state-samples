# Worker sample 

worker process

```bash
$ docker-compose up -d --build

# prepare state
$ docker-compose exec goapp curl -X POST localhost:3500/v1.0/state/redis-state -H 'Content-Type: "application/json"' -d '[{"key": "m-1", "value": {"id": "m-1", "message": "hello!"}}]'
$ docker-compose exec redis redis-cli hgetall "goapp||m-1"
1) "data"
2) "{\"id\":\"m-1\",\"message\":\"hello!\"}"
3) "version"
4) "1"

# publish appending message event to worker
$ docker-compose exec goapp curl -X POST localhost:3000/append-message -H 'Content-Type: "application/json"' -d '{"event_id": "e-1", "message_id": "m-1", "message": " gopher"}'

# (in another window) message updated gradually
$ docker-compose exec redis redis-cli hgetall "goapp||m-1"
1) "data"
2) "{\"id\":\"m-1\",\"message\":\"hello! \"}"
3) "version"
4) "2"
$ docker-compose exec redis redis-cli hgetall "goapp||m-1"
1) "data"
2) "{\"id\":\"m-1\",\"message\":\"hello! g\"}"
3) "version"
4) "3"
# ...
$ docker-compose exec redis redis-cli hgetall "goapp||m-1"
1) "data"
2) "{\"id\":\"m-1\",\"message\":\"hello! gopher\"}"
3) "version"
4) "8"
```

workers conflict

```bash
$ docker-compose up -d --build
$ docker-compose exec goapp curl -X POST localhost:3500/v1.0/state/redis-state -H 'Content-Type: "application/json"' -d '[{"key": "m-1", "value": {"id": "m-1", "message": "hello!"}}]'

# publish appending 2 events to each worker at the same time
$ docker-compose exec goapp curl -X POST localhost:3000/append-message -H 'Content-Type: "application/json"' -d '{"event_id": "e-1", "message_id": "m-1", "message": " gopher"}' & docker-compose exec goapp2 curl -X POST localhost:3000/append-message -H 'Content-Type: "application/json"' -d '{"event_id": "e-1", "message_id": "m-1", "message": " gopher"}' &

# (in another window) message updated but entangled
$ docker-compose exec redis redis-cli hgetall "goapp||m-1"
1) "data"
2) "{\"id\":\"m-1\",\"message\":\"hello!\"}"
3) "version"
4) "1"
$ docker-compose exec redis redis-cli hgetall "goapp||m-1"
1) "data"
2) "{\"id\":\"m-1\",\"message\":\"hello! g\"}"
3) "version"
4) "4"
$ docker-compose exec redis redis-cli hgetall "goapp||m-1"
1) "data"
2) "{\"id\":\"m-1\",\"message\":\"hello! ggp\"}"
3) "version"
4) "7"
# ...
$ docker-compose exec redis redis-cli hgetall "goapp||m-1"
1) "data"
2) "{\"id\":\"m-1\",\"message\":\"hello! ggpoerher\"}"
3) "version"
4) "15"
```
