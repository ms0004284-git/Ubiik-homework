package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"message-gateway-service/models"
)

func bindMessageData(input interface{}, output interface{}) error {
	bytes, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, output)
}

func buildReadingRequest(message models.Message, device models.Device, readData models.ReadingData, readReq *models.ReadingRequest) error {
	if message.DeviceID == "" {
		return errors.New("DeviceID is missing")
	}
	if device.Username == "" {
		return errors.New("Username is missing")
	}
	if readData.Reading == nil {
		return errors.New("Reading data is missing")
	}

	readReq.DeviceID = message.DeviceID
	readReq.Username = device.Username
	readReq.Reading = readData.Reading

	return nil
}

func HandleMessage(ctx *gin.Context) {
	var message models.Message
	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := resty.New()
	switch message.Type {
	case "registration":
		handleRegistration(ctx, client, message)
	case "reading":
		handleReading(ctx, client, message)
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unknown message type."})
	}
}

func handleRegistration(ctx *gin.Context, client *resty.Client, message models.Message) {
	var regData models.RegistrationData
	if err := bindMessageData(message.Data, &regData); err != nil || regData.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration data."})
		return
	}

	url := fmt.Sprintf("http://device-service:8080/devices/%s", message.DeviceID)
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(regData).
		Put(url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to device-service."})
	}
}

func handleReading(ctx *gin.Context, client *resty.Client, message models.Message) {
	var readData models.ReadingData
	if err := bindMessageData(message.Data, &readData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reading data"})
		return
	}

	url := fmt.Sprintf("http://device-service:8080/devices/%s/username", message.DeviceID)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		Get(url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to device-service"})
		return
	}

	var device models.Device
	if err := json.Unmarshal(resp.Body(), &device); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse device data"})
		return
	}

	var readReq models.ReadingRequest
	if err := buildReadingRequest(message, device, readData, &readReq); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	url = "http://reading-service/readings"
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(readReq).
		Post(url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to reading-service"})
	}
}
