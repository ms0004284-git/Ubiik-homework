
## File Description
```
├─ Dockerfile
├─ REQUIREMENT.md
├─ README.md
├─ docker-compose.yml
└─ iotservice
   ├─ device-service
   │  ├─ Dockerfile
   │  ├─ db
   │  │  └─ db.go
   │  ├─ go.mod
   │  ├─ go.sum
   │  └─ main.go
   └─ message-gateway-service
      ├─ Dockerfile
      ├─ go.mod
      ├─ go.sum
      └─ main.go

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

## Features

1. **device-service**
   - Provides endpoints to create, update, and retrieve mappings between `deviceId` and `username`.

2. **message-gateway-service**
   - Processes incoming IoT messages.
   - Integrates with `device-service` to handle device registration and readings.
