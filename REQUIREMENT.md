# Simple IoT Data Processing System

Design and implement two HTTP microservices for a simple IoT data processing system.

## **Service 1: device-service**

This service is responsible for managing the mapping between `deviceId` and `username`.

### **Endpoints**

1. **PUT /devices/:deviceId**

    - **Request Body**:
        ```json
        {
            "username": "my-username"
        }
        ```

    - **Description**:
        - `deviceId` and `username` can be any string.
        - If the `deviceId` already exists, update the associated `username`.
        - Otherwise, create a new mapping for the `deviceId` and `username`.

2. **GET /devices/:deviceId/username**

    - **Response Body**:
        ```json
        {
            "username": "my-username"
        }
        ```

    - **Description**:
        - Retrieve the `username` associated with the given `deviceId`.

## **Service 2: message-gateway-service**

This service handles incoming messages from IoT devices based on the message type.

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

## **Requirements and Constraints**

- **Programming Language**:
    - You may use any programming language of your choice.

- **Data Persistence**:
    - You may use any data storage solution of your choice.

- **Flexibility**:
    - You can design and implement any additional details by yourself (e.g., data structures, error handling).

- **Additional Notes**:
    - Assume `reading-service` (another existing microservice) mentioned above is developed by another team, so you just need to follow the API format to send requests to it.

## **Evaluation Criteria**

- **Functionality**:
    - Correctness of API behavior.
    - Handling of edge cases.

- **Code Quality**:
    - Clean and well-structured code.

- **Bonus Points**:
    - Unit tests.
    - Instructions to run the services locally (e.g., Dockerfile, README).
    - Consider how to handle concurrent requests to ensure data consistency.
