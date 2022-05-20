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

func GetAuthorizationToken(authJson string) (AuthorizationToken, error) {
	reader := strings.NewReader(authJson)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/oauth/authorize", HostURL), reader)
	if err != nil {
		return AuthorizationToken{}, err
	}

	c := http.Client{Timeout: 5 * time.Second}
	c = *httptrace.WrapClient(&c)
	res, err := c.Do(req)
	if err != nil {
		return AuthorizationToken{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return AuthorizationToken{}, err
	}

	if res.StatusCode != http.StatusOK {
		return AuthorizationToken{}, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	var authToken AuthorizationToken
	err = json.Unmarshal(body, &authToken)
	if err != nil {
		return AuthorizationToken{}, err
	}

	return authToken, nil
}

func GetAccessToken(authToken string) (AccessToken, error) {
	requestObj := map[string]interface{}{
		"authorization": authToken,
	}
	requestBody, err := json.Marshal(requestObj)

	reader := strings.NewReader(string(requestBody))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/oauth/accesstoken", HostURL), reader)
	if err != nil {
		return AccessToken{}, err
	}

	c := http.Client{Timeout: 5 * time.Second}
	c = *httptrace.WrapClient(&c)
	res, err := c.Do(req)
	if err != nil {
		return AccessToken{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return AccessToken{}, err
	}

	if res.StatusCode != http.StatusOK {
		return AccessToken{}, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	var accessToken AccessToken
	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		return AccessToken{}, err
	}

	return accessToken, nil
}
