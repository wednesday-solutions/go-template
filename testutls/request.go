package testutls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

type RequestParameters struct {
	Pathname    string
	E           *echo.Echo
	RequestBody string
	AuthHeader  string
	HttpMethod  string
	IsGraphQL   bool
}

func MakeRequest(parameters RequestParameters) (map[string]interface{}, error) {
	_, _, jsonRes, err := MakeAndGetRequest(parameters)
	return jsonRes, err
}

func SimpleMakeRequest(parameters RequestParameters) (*http.Request, *http.Response, error) {
	client := &http.Client{}
	ts := httptest.NewServer(parameters.E)
	path := ts.URL + parameters.Pathname
	defer ts.Close()

	req, _ := http.NewRequest(
		parameters.HttpMethod,
		path,
		bytes.NewBuffer([]byte(parameters.RequestBody)),
	)

	req.Header.Set("authorization", parameters.AuthHeader)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
		log.Fatal("Cannot create http request")
	}
	return req, res, nil
}

func MakeAndGetRequest(parameters RequestParameters) (*http.Request, *http.Response, map[string]interface{}, error) {
	req, res, err := SimpleMakeRequest(parameters)
	if err != nil {
		return nil, nil, nil, err
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, nil, err
	}
	if parameters.IsGraphQL {
		var jsonRes map[string]interface{}
		err = json.Unmarshal(bodyBytes, &jsonRes)
		fmt.Print(err, jsonRes)
		if err != nil {
			return nil, nil, nil, err

		}
		return req, res, jsonRes, nil
	} else {
		var jsonRes map[string]interface{}
		err = json.Unmarshal(bodyBytes, &jsonRes)
		if err != nil {
			return nil, nil, nil, err
		}
		return req, res, jsonRes, nil
	}
}
