package sensorpush

import "time"

type AccessToken struct {
	Token string `json:"accesstoken"`
}

type AuthorizationToken struct {
	Token  string `json:"authorization"`
	Apikey string `json:"apikey"`
}

type SamplesResponse struct {
	LastTime time.Time                  `json:"last_time"`
	Sensors  map[string][]SensorReading `json:"sensors"`
}

type Sensor struct {
	Id             string  `json:"id"`
	Name           string  `json:"name"`
	Active         bool    `json:"active"`
	Address        string  `json:"address"`
	DeviceId       string  `json:"deviceId"`
	Rssi           int     `json:"rssi"`
	Type           string  `json:"type"`
	BatteryVoltage float64 `json:"battery_voltage"`
}

type SensorReading struct {
	Timestamp   time.Time `json:"observed"`
	Gateway     string    `json:"gateways"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
}

type Gateway struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	LastSeen time.Time `json:"last_seen"`
	Version  string    `json:"version"`
}
