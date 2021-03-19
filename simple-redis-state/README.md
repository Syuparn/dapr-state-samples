# Simple Redis state sample 

```bash
$ docker-compose up -d --build
# request to goapp
$ curl -X POST localhost:8080/messages -H "Content-Type: application/json" -d '{"message": "hello, state!"}'
{"id":"01F15A29CH1TS0QZG0NKN8BMRQ","message":"hello, state!"}

# the message saved in redis
$ docker-compose exec redis redis-cli keys '*'
1) "goapp||01F15A29CH1TS0QZG0NKN8BMRQ"
$ docker-compose exec redis redis-cli hgetall "goapp||01F15A29CH1TS0QZG0NKN8BMRQ"
1) "data"
2) "{\"id\":\"01F15A29CH1TS0QZG0NKN8BMRQ\",\"message\":\"hello, state!\"}"
3) "version"
4) "1"
```