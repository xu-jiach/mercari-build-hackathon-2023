## Simple mercari API

### Run server

```shell
$ cd backend # move to mercari-build-hackathon-2023/backend
$ go run main.go
```

Please call this endpoint for initialize data. 

```shell
$ curl -X POST 'http://127.0.0.1:9000/initialize'
```


### Spec

| Features                           | Endpoint                         | Benchmarker spec                                                                                                        |
|------------------------------------|----------------------------------|-------------------------------------------------------------------------------------------------------------------------|
| Reset db for bench                 | `POST /initialize`               | This endpoint will be called before bench. <br>The endpoint reset database data. <br>The endpoint have to finish 10 sec |
| Access log                         | `GET /log`                       | Show access log. This endpoint is not target of scoring. Check after bench and change freely.                           |
| User Registration                  | `POST /register`                 |                                                                                                                         |
| Login                              | `POST /login`                    |                                                                                                                         |
| List of items                      | `GET /items`                     | The benchmarker ensures that at least 12 items are returned if exist.                                                   |
| Item detail                        | `GET /items/:itemID`             |                                                                                                                         |
| Item image                         | `GET /items/:itemID/image`       | Don't change image. Benchmarker will send images up to 1MB in size.                                                     |
| Search item by name *unimplemented | `GET /search?name=<search word>` | Response item have to Include search word <br>The benchmarker ensures that at least 12 items are returned if exist.     |
| Get balance                        | `GET /balance`                   |                                                                                                                         |
| Add balance                        | `POST /balance`                  |                                                                                                                         |
| User listed item                   | `/users/:userID/items`           | Sort by created time                                                                                                    |
| Item detail                        | `GET /items/:itemID`             |                                                                                                                         |
| Purchase item                      | `POST /purchase/:itemID`         |                                                                                                                         |
| Edit item *unimplemented           | `PUT /items `                    | Expect same request body as POST /items                                                                                 |
| Create new item draft              | `POST /items`                    |                                                                                                                         |
| Start to sell item                 | `POST /sell`                     |                                                                                                                         |


### Backend scoring
The Backend API will be evaluated by a benchmark tester.  
The benchmark tester will conduct tests on the endpoints specified in the Spec.
You can run bench from `Run Benchmark` button on [dashboard](https://mercari-build-hackathon-2023-front-d3sqdyhc4a-uc.a.run.app/).
It is possible to add to existing request-response cycles, but be careful not to reduce the endpoints, as this will cause the benchmark tester's validation to fail.  
Don't call put endpoint while benchmarker running when you start bench.

**Scoring**

* Success to POST request in 1 second ... 3pt
* Success to GET request in 1 second ... 1pt
  * The following endpoints are exception
    * `/search` ... 5pt
* Validation was failed ... -10pt
    * examples
        * When calling `/purchase/1` 3 times, the seller's balance was added 3 times.
        * User balance is under 0yen.
* Benchmarker will be stopped.
    * Failed to initialize
    * More than 10 calls timed out.
    * Too many invalid response.
* If a response is expected to have a status code of 200 but actually differs, the benchmark execution time will be reduced by 0.3 seconds per one difference.

### Example
Sample code for calling some endpoints after running server

サーバを立ち上げた後、curlでテストする場合のサンプル

```shell
# Registration
# {"id":11,"name":"momom"}
$ curl -X POST 'http://127.0.0.1:9000/register' -d '{"name": "momom", "password": "password"}'  -H 'Content-Type: application/json'
# Login (get login token)
# {"id":11,"name":"momom","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMSwiZXhwIjoxNjg0NTgxNjU3fQ.7YGvgOsKI1EIr8a9yw0Ny6GRmmUJjrAkjjypdpj74qw"}
$ curl -i -X POST 'http://127.0.0.1:9000/login' -d '{"user_id": 11, "password": "password"}'  -H 'Content-Type: application/json'
# Add item
# Please put image.jpg on backend folder to call this endpoint 
# {"id":21}
$ curl -X POST \
  --url 'http://127.0.0.1:9000/items' \
  -F 'name=item' \
  -F 'category_id=1' \
  -F 'price=100' \
  -F 'description=samplesamplesample' \
  -F 'image=@image.jpg' \
  -H "Authorization: Bearer <Token which get login endpoint>"
# Item list
# [{"id":3,"name":"Cucumber","price":80,"image": ..."}]
curl -X GET 'http://127.0.0.1:9000/users/1/items' -H "Authorization: Bearer <ログイン時のレスポンスで返ってきたtokenの値を入れる>"
# Add a balance 
# "successful"
curl -X POST 'http://127.0.0.1:9000/balance' -d '{"balance": 1000}' -H "Authorization: Bearer <ログイン時のレスポンスで返ってきたtokenの値を入れる>" -H 'Content-Type: application/json'
# See a balance
# {"balance":1000}
curl -X GET 'http://127.0.0.1:9000/balance' -H "Authorization: Bearer <ログイン時のレスポンスで返ってきたtokenの値を入れる>"
# Sell
# "successful"
curl -X POST 'http://127.0.0.1:9000/sell' -d '{"user_id": 1, "item_id": 1}' -H "Authorization: Bearer <ログイン時のレスポンスで返ってきたtokenの値を入れる>" -H 'Content-Type: application/json'
# Purchase
# "successful"
curl -X POST 'http://127.0.0.1:9000/purchase/1' -H "Authorization: Bearer <ログイン時のレスポンスで返ってきたtokenの値を入れる>" -H 'Content-Type: application/json'
```

###  Structure

```
backend
├── README.md
├── db # データベース関連のソースコード
│   ├── driver.go
│   ├── repository.go
│   └── schema.sql
├── domain # アイテム、ユーザにまつわる定義群
│   ├── item.go
│   └── user.go
├── handler # 各エンドポイントのソースコード
│   └── handler.go
├── main.go # サーバの立ち上げ/シャットダウンのソースコード
└── swagger.yaml # APIドキュメント
```
