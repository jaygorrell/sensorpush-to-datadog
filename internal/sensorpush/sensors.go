package sensorpush

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

func GetSensors(accessToken string) (map[string]Sensor, error) {
	requestObj := map[string]interface{}{
		"active": true,
	}
	requestBody, err := json.Marshal(requestObj)

	reader := strings.NewReader(string(requestBody))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/devices/sensors", HostURL), reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("authorization", accessToken)

	c := http.Client{Timeout: 5 * time.Second}
	c = *httptrace.WrapClient(&c)
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	var sensors map[string]Sensor
	err = json.Unmarshal(body, &sensors)
	if err != nil {
		return nil, err
	}

	return sensors, nil
}

func GetReadings(accessToken, sensorId string) ([]SensorReading, error) {
	requestObj := map[string]interface{}{
		"active":    true,
		"limit":     1,
		"sensors":   []string{sensorId},
		"startTime": time.Now().Add(time.Duration(-3) * time.Minute),
	}
	requestBody, err := json.Marshal(requestObj)

	reader := strings.NewReader(string(requestBody))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/samples", HostURL), reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("authorization", accessToken)

	c := http.Client{Timeout: 5 * time.Second}
	c = *httptrace.WrapClient(&c)
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	var resObj SamplesResponse
	err = json.Unmarshal(body, &resObj)
	if err != nil {
		return nil, err
	}

	var readings []SensorReading

	for _, reading := range resObj.Sensors[sensorId] {
		readings = append(readings, SensorReading{
			Timestamp:   reading.Timestamp,
			Gateway:     reading.Gateway,
			Temperature: reading.Temperature,
			Humidity:    reading.Humidity,
		})
	}

	return readings, nil
}
