
## File Description
```
├─ .gitignore
├─ Dockerfile
├─ README.md
├─ REQUIREMENT.md
├─ docker-compose.yml
└─ iotservice
   ├─ device-service
   │  ├─ Dockerfile
   │  ├─ db
   │  │  └─ db.go
   │  ├─ go.mod
   │  ├─ go.sum
   │  ├─ handlers
   │  │  ├─ handler.go
   │  │  └─ handler_test.go
   │  ├─ main.go
   │  ├─ models
   │  │  └─ models.go
   │  └─ repositories
   │     └─ repo.go
   └─ message-gateway-service
      ├─ Dockerfile
      ├─ go.mod
      ├─ go.sum
      ├─ handlers
      │  ├─ handler.go
      │  └─ handler_test.go
      ├─ main.go
      └─ models
         └─ models.go

```

## Quick Start
透過 docker-compose 建置 image 並開啟服務
```
docker-compose build
docker-compose up -d 
```

可以透過 go-dev container 來進入環境開發
```
docker exec -it go-dev bash
```

檢視資料庫狀態
```
docker exec -it mysql bash
mysql -u root -p
USE iot
```

進入 device-service or message-gateway-service container 中進行 unit test
```
docker exec -it <service> bash // <service> is device-service or message-gateway-service.
go test ./handlers 
```

## Features

1. **device-service**
   - Provides endpoints to create, update, and retrieve mappings between `deviceId` and `username`.

2. **message-gateway-service**
   - Processes incoming IoT messages.
   - Integrates with `device-service` to handle device registration and readings.
