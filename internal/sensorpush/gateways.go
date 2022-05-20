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

func GetGateways(accessToken string) (map[string]Gateway, error) {
	requestObj := map[string]interface{}{}
	requestBody, err := json.Marshal(requestObj)

	reader := strings.NewReader(string(requestBody))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/devices/gateways", HostURL), reader)
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

	var gateways map[string]Gateway
	err = json.Unmarshal(body, &gateways)
	if err != nil {
		return nil, err
	}

	return gateways, nil
}
