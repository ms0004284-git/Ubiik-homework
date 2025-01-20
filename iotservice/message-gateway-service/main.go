package main

import (
	"encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"errors"
	"fmt"
)

type Device struct {
	Username	string		`json:"username"`
}

type Message struct {
	DeviceID   	string    	`json:"deviceId"`
	Type 		string 		`json:"type"`
	Data   		interface{} `json:"data"`
}

type RegistrationData struct {
	Username 	string 		`json:"username"`
}

type ReadingData struct {
	Reading 	float64		`json:"reading"`
}

type ReadingRequest struct {
	DeviceID   	string    	`json:"deviceId"`
	Username	string 		`json:"username"`
	Reading 	float64 	`json:"reading"`
}


func bindMessageData(input interface{}, output interface{}) error {
	bytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, output)
}

func buildReadingRequest(message Message, device Device, readData ReadingData, readReq *ReadingRequest) error {
	if message.DeviceID == "" {
		return errors.New("DeviceID is missing in the message")
	}
	if device.Username == "" {
		return errors.New("Username is missing in the device")
	}
	if readData.Reading <= 0 {
		return errors.New("Reading value must be greater than 0")
	}

	readReq.DeviceID = message.DeviceID
	readReq.Username = device.Username
	readReq.Reading = readData.Reading

	return nil
}

func main() {
	r := gin.Default()

	r.POST("/messages", func(ctx *gin.Context){
		var message Message
		if err := ctx.ShouldBindJSON(&message); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		client := resty.New()
		switch message.Type {
			case "registration":
				// binding Message.Data
				var regData RegistrationData
				if err := bindMessageData(message.Data, &regData); err != nil || regData.Username == "" {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration data"})
					return
				}

				// send request to device-service
				url := fmt.Sprintf("http://device-service:8080/devices/%s", message.DeviceID)
				_, err := client.R().
					SetHeader("Content-Type", "application/json").
					SetBody(regData).
					Put(url) 
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to device-service"})
					return
				}
				
			case "reading":
				// binding Message.Data
				var readData ReadingData
				if err := bindMessageData(message.Data, &readData); err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reading data"})
					return
				}
				fmt.Print(readData.Reading)

				// send request to device-service
				url := fmt.Sprintf("http://device-service:8080/devices/%s/username", message.DeviceID)
				resp, err := client.R().
					SetHeader("Content-Type", "application/json").
					SetBody(readData).
					Get(url) 
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to device-service"})
					return
				}
				
				// binding Device
				var device Device
				if err := json.Unmarshal(resp.Body(), &device); err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse device data"})
					return
				}
				
				// binding ReadingRequest
				var readReq ReadingRequest
				if err := buildReadingRequest(message, device, readData, &readReq); err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				// send request to reading-service
				url = fmt.Sprintf("http://reading-service/readings") // mock url
				resp, err = client.R().
					SetHeader("Content-Type", "application/json").
					SetBody(readReq).
					Post(url)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to reading-service"})
					return
				}

			default:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unknown message type"})
		}
		

	})

	r.Run(":8081") 
}