package models

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