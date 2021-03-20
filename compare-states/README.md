# Simple Redis state sample 

```bash
$ docker-compose up -d --build
# request to goapp
$ curl -X POST localhost:8080/messages -H "Content-Type: application/json" -d '{"message": "hello, state!"}'
{"id":"01F1760YWJDABNNZ2QRSFTQH6J","message":"hello, state!"}

# the message saved in redis
$ docker-compose exec redis redis-cli keys '*'
1) "goapp||01F1760YWJDABNNZ2QRSFTQH6J"
$ docker-compose exec redis redis-cli hgetall "goapp||01F1760YWJDABNNZ2QRSFTQH6J"
1) "data"
2) "{\"id\":\"01F1760YWJDABNNZ2QRSFTQH6J\",\"message\":\"hello, state!\"}"
3) "version"
4) "1"

# the message saved in mysql
$ docker-compose exec mysql mysql -u root -p
mysql> use dapr_state_store;
mysql> SELECT * FROM state\G
*************************** 1. row ***************************
        id: goapp||01F1760YWJDABNNZ2QRSFTQH6J
     value: {"id": "01F1760YWJDABNNZ2QRSFTQH6J", "message": "hello, state!"}
insertDate: 2021-03-20 06:44:09
updateDate: 2021-03-20 06:44:09
      eTag: e6dc50aa-214f-4eaf-9818-022bc4bcd28f
1 row in set (0.00 sec)
```
