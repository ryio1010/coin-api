# G-CodingTest

## 技術スタック

- Go
- Gin
- PostgreSQL
- Docker
- Gorm
- zerolog
- CleanArchitecture

## 起動方法

クローン
1. git clone git@github.com:ryio1010/coin-api.git

Dockerの起動
1. docker-compose build
2. docker-compose up -d

## API実行方法

- ユーザー登録
    - method : POST
    - URL : localhost:8081/v1/user
    - RequestJsonBody : {"username":"test1","password":"test1"}

- 対象ユーザー残高取得
    - method : GET
    - URL : localhost:8081/v1/user/{userid}
    - RequestJsonBody : なし

- コイン履歴確認
    - method : GET
    - URL : localhost:8081/v1/coin/{userid}
    - RequestJsonBody : なし

- コイン追加
    - method : PUT
    - URL : localhost:8081/v1/coin
    - RequestJsonBody : {"userid": "1","operation": "ADD","amount": "1000"}

- コイン消費
    - method : PUT
    - URL : localhost:8081/v1/coin
    - RequestJsonBody : {"userid": "1","operation": "USE","amount": "100"}

- コイン送金
    - method : PUT
    - URL : localhost:8081/v1/coin/send
    - RequestJsonBody : {"sender": "1","receiver": "2","amount": "100"}
