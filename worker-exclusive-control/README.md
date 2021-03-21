# Exclusive control sample 

```bash
$ docker-compose up -d --build
$ docker-compose exec goapp curl -X POST localhost:3500/v1.0/state/redis-state -H 'Content-Type: "application/json"' -d '[{"key": "m-1", "value": {"id": "m-1", "message": "hello!"}}]'

# publish appending 2 events to each worker at the same time
$ docker-compose exec goapp curl -X POST localhost:3000/append-message -H 'Content-Type: "application/json"' -d '{"event_id": "e-1", "message_id": "m-1", "message": " gopher"}' & docker-compose exec goapp2 curl -X POST localhost:3000/append-message -H 'Content-Type: "application/json"' -d '{"event_id": "e-1", "message_id": "m-1", "message": " gopher"}' &

# (in another window) message updated without conflict
$ docker-compose exec redis redis-cli hgetall "goapp||m-1"
1) "data"
2) "{\"id\":\"m-1\",\"message\":\"hello!\"}"
3) "version"
4) "1"
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

# log shows conflicted worker failed to save
$ docker-compose logs --tail 1 -f goapp goapp2 
Attaching to worker-exclusive-control_goapp_1, worker-exclusive-control_goapp2_1
goapp2_1       | dapr client initializing for: 127.0.0.1:50001
goapp_1        | dapr client initializing for: 127.0.0.1:50001
goapp2_1       | 2021/03/21 01:54:55 binding - Data:{"event_id": "e-1", "message_id": "m-1", "message": " gopher"}
goapp2_1       | 2021/03/21 01:54:55 event: {e-1 m-1  gopher}
goapp2_1       | 2021/03/21 01:54:55 message: {m-1 hello!}
goapp_1        | 2021/03/21 01:54:55 binding - Data:{"event_id": "e-1", "message_id": "m-1", "message": " gopher"}
goapp_1        | 2021/03/21 01:54:55 event: {e-1 m-1  gopher}
goapp_1        | 2021/03/21 01:54:55 message: {m-1 hello!}
goapp_1        | 2021/03/21 01:54:56 updated message: {m-1 hello! }
goapp2_1       | 2021/03/21 01:54:56 updated message: {m-1 hello! }
goapp2_1       | 2021/03/21 01:54:56 failed to save message: {m-1 hello! }
goapp_1        | 2021/03/21 01:54:58 message: {m-1 hello! }
goapp_1        | 2021/03/21 01:55:00 updated message: {m-1 hello! g}
goapp_1        | 2021/03/21 01:55:02 message: {m-1 hello! g}
goapp_1        | 2021/03/21 01:55:03 updated message: {m-1 hello! go}
goapp_1        | 2021/03/21 01:55:05 message: {m-1 hello! go}
goapp_1        | 2021/03/21 01:55:05 updated message: {m-1 hello! gop}
goapp_1        | 2021/03/21 01:55:07 message: {m-1 hello! gop}
goapp_1        | 2021/03/21 01:55:09 updated message: {m-1 hello! goph}
goapp_1        | 2021/03/21 01:55:10 message: {m-1 hello! goph}
goapp_1        | 2021/03/21 01:55:10 updated message: {m-1 hello! gophe}
goapp_1        | 2021/03/21 01:55:11 message: {m-1 hello! gophe}
goapp_1        | 2021/03/21 01:55:13 updated message: {m-1 hello! gopher}
```
