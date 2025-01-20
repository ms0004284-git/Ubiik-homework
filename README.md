
# File Description
```
ann
├─ .DS_Store
├─ Dockerfile
├─ README 1.md
├─ docker-compose.yml
├─ homework
│  ├─ device-service
│  │  ├─ Dockerfile
│  │  ├─ db
│  │  │  └─ db.go
│  │  ├─ go.mod
│  │  ├─ go.sum
│  │  └─ main.go
│  └─ message-gateway-service
│     ├─ Dockerfile
│     ├─ go.mod
│     ├─ go.sum
│     └─ main.go
└─ testcase.txt

```

# Quick Start
透過 docker-compose 建置 image 並開啟服務
```
docker-compose build
docker-compose up -d 
```

檢視資料庫狀態
```
docker exec -it mysql bash
mysql -u root -p
USE iot
```

# Requirement

### **Endpoints**

1. **POST /messages**

    - **Request Body**:
        ```json
        {
            "deviceId": "my-device",
            "type": "registration",
            "data": {
                "username": "my-username"
            }
        }
        ```

    - **Description**:
        - Handle incoming messages from IoT devices.
        - Fields:
            - `deviceId`: can be any string.
            - `type`: can be `"registration"` or `"reading"`.
            - `data`: can be different format depends on `type`.
                - When `type` is `"registration"`, the `data` will be in this format:
                    ```json
                    {
                        "username": "my-username"
                    }
                    ```
                - When `type` is `"reading"`, the `data` will be in this format:
                    ```json
                    {
                        "reading": 20.5
                    }
                    ```
        - Behavior:
            - When `type` is `"registration"`:
                - Send a `PUT /devices/:deviceId` request to `device-service` (Service 1) to upsert the username of `deviceId`.
            - When `type` is `"reading"`:
                - Send a `GET /devices/:deviceId/username` request to `device-service` to retrieve the `username`.
                - Send a `POST /readings` request to `reading-service` (another existing microservice) to store the reading:
                    ```json
                    {
                        "deviceId": "my-device",
                        "username": "my-username",
                        "reading": 20.5
                    }
                    ```
