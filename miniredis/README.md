# Miniredis state sample 

Simulate abnormal scenario of Redis state by Miniredis

```bash
$ docker-compose up -d --build
# request to goapp
$ curl -X POST localhost:8080/messages -H "Content-Type: application/json" -d '{"message": "hello, state!"}'
{"id":"01F15A29CH1TS0QZG0NKN8BMRQ","message":"hello, state!"}

# But... MiniRedis terminates after 10th request!!
$ curl -X POST localhost:8080/messages -H "Content-Type: application/json" -d '{"message": "hello, state!"}'
{"message":"Internal Server Error"}

$ docker-compose logs redis
Attaching to miniredis_redis_1
redis_1       | 2021/04/11 02:17:34 miniredis serves on [::]:6379
redis_1       | 2021/04/11 02:17:36 Request: 1/10
redis_1       | 2021/04/11 02:17:36 Request: 2/10
...
redis_1       | 2021/04/11 02:24:46 Request: 10/10
redis_1       | 2021/04/11 02:24:46 Bye!
$ docker-compose logs --tail 1 goapp
Attaching to miniredis_goapp_1
goapp_1       | {"time":"2021-04-11T02:27:46.768353137Z","id":"","remote_ip":"172.27.0.1","host":"localhost:8080","method":"POST","uri":"/messages","user_agent":"curl/7.68.0","status":500,"error":"failed to save message: failed to save message: error saving state: rpc error: code = Internal desc = failed saving state in state store redis-state: failed to set key goapp||01F2ZC3AYERA1CYBY6XS612ZM5: MiniRedis went home.","latency":1422117,"latency_human":"1.422117ms","bytes_in":28,"bytes_out":36}
```
