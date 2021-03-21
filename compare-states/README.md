# State comparison sample 

```bash
$ docker-compose up -d --build
# request to goapp
$ curl -X POST localhost:8080/messages -H "Content-Type: application/json" -d '{"message": "hello, state!"}'
{"id":"01F178XVG5HS2W3658TKQA6DA1","message":"hello, state!"}

# the message saved in redis
$ docker-compose exec redis redis-cli keys '*'
1) "goapp||01F178XVG5HS2W3658TKQA6DA1"
$ docker-compose exec redis redis-cli hgetall "goapp||01F178XVG5HS2W3658TKQA6DA1"
1) "data"
2) "{\"id\":\"01F178XVG5HS2W3658TKQA6DA1\",\"message\":\"hello, state!\"}"
3) "version"
4) "1"

# the message saved in mysql
$ docker-compose exec mysql mysql -u root -p
mysql> use dapr_state_store;
mysql> SELECT * FROM state\G
*************************** 1. row ***************************
        id: goapp||01F178XVG5HS2W3658TKQA6DA1
     value: {"id": "01F178XVG5HS2W3658TKQA6DA1", "message": "hello, state!"}
insertDate: 2021-03-20 07:34:53
updateDate: 2021-03-20 07:34:53
      eTag: 2ba1c054-4b28-4847-acd0-cf25ed542732
1 rows in set (0.00 sec)

# the message saved in mongodb
$ docker-compose exec mongodb mongo
> use daprStore;
switched to db daprStore
> db.daprCollection.find()
{ "_id" : "goapp||01F178XVG5HS2W3658TKQA6DA1", "_etag" : "9ad9369d-d6b3-4df1-b94e-282abf344c7a", "value" : "{\"id\":\"01F178XVG5HS2W3658TKQA6DA1\",\"message\":\"hello, state!\"}" }
```
