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

	if err := PutDeviceService(ctx, client, message.DeviceID, regData); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to device-service."})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success to register."})
}

func handleReading(ctx *gin.Context, client *resty.Client, message models.Message) {
	var readData models.ReadingData
	if err := bindMessageData(message.Data, &readData); err != nil || readData.Reading == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reading data."})
		return
	}

	resp ,err := GetDeviceService(ctx, client, message.DeviceID)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to device-service"})
		return
	}

	var device models.Device
	if err := json.Unmarshal(resp.Body(), &device); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse device data"})
		return
	}

	var readReq models.ReadingRequest
	if err := buildReadingRequest(message, device, readData, &readReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := PostReadingService(ctx, client, readReq); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to reading-service."})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success to send request to reading-service."})
}

func PostReadingService(ctx *gin.Context, client *resty.Client, readReq models.ReadingRequest) error{
	url := "http://reading-service/readings"
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(readReq).
		Post(url)

	if err != nil {
		return err
	}
	return nil
}

func GetDeviceService(ctx *gin.Context, client *resty.Client, deviceID string) (*resty.Response, error){
	url := fmt.Sprintf("http://device-service:8080/devices/%s/username", deviceID)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func PutDeviceService(ctx *gin.Context, client *resty.Client, deviceID string, regData models.RegistrationData) error{
	url := fmt.Sprintf("http://device-service:8080/devices/%s", deviceID)
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(regData).
		Put(url)
	if err != nil {
		return err
	}
	return nil
}



