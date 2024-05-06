# curl's

curl --verbose 127.0.0.1:4221/echo/pineapple
```
REQUEST
GET /echo/pineapple HTTP/1.1
Host: localhost:4221
User-Agent: curl/7.64.1
```

```
RESPONSE
HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 9

pineapple
```

---  
curl --verbose 127.0.0.1:4221/user-agent
```
REQUEST
GET /user-agent HTTP/1.1
Host: localhost:4221
User-Agent: curl/7.64.1
```
```
RESPONSE
HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 11

curl/7.64.1
```

---

curl -v -X GET http://localhost:4221/files/mango_strawberry_strawberry_pineapple

---

curl -v -X POST http://localhost:4221/files/raspberry_pineapple_strawberry_strawberry -d 'pear pineapple grape raspberry apple pear grape strawberry'
