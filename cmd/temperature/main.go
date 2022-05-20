package main

import (
	"context"
	"fmt"
	"os"

	"github.com/DataDog/datadog-lambda-go"
	"github.com/aws/aws-lambda-go/lambda"

	"home-sensor-cron/internal/aws"
	"home-sensor-cron/internal/sensorpush"
)

func main() {
	lambda.Start(ddlambda.WrapFunction(handleTemperature, nil))
}

func handleTemperature(ctx context.Context) error {
	secretName := os.Getenv("AWS_SENSORPUSH_SECRET")
	awsRegion := os.Getenv("REGION")

	authJson, err := aws.GetSecret(secretName, awsRegion)
	if err != nil {
		fmt.Println(err)
		return err
	}

	authToken, err := sensorpush.GetAuthorizationToken(authJson)
	if err != nil {
		fmt.Println(err)
		return err
	}

	accessToken, err := sensorpush.GetAccessToken(authToken.Token)
	if err != nil {
		fmt.Println(err)
		return err
	}

	gateways, err := sensorpush.GetGateways(accessToken.Token)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// TODO: map gateway id to name for use in sensor values
	for _, gateway := range gateways {
		ddlambda.MetricWithTimestamp(
			"sensorpush.gateway.running",
			1,
			gateway.LastSeen,
			fmt.Sprintf("name:%s", gateway.Name),
			fmt.Sprintf("id:%s", gateway.Id),
			fmt.Sprintf("version:%s", gateway.Version),
			fmt.Sprintf("env:%s", os.Getenv("ENV")),
		)
	}

	sensors, err := sensorpush.GetSensors(accessToken.Token)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, sensor := range sensors {
		ddlambda.Metric(
			"sensorpush.sensor.running",
			1,
			fmt.Sprintf("name:%s", sensor.Name),
			fmt.Sprintf("type:%s", sensor.Type),
			fmt.Sprintf("env:%s", os.Getenv("ENV")),
		)

		readings, err := sensorpush.GetReadings(accessToken.Token, sensor.Id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		for _, reading := range readings {
			fmt.Println(reading)
			ddlambda.MetricWithTimestamp(
				"sensorpush.sensor.temperature",
				reading.Temperature,
				reading.Timestamp,
				fmt.Sprintf("name:%s", sensor.Name),
				fmt.Sprintf("type:%s", sensor.Type),
				fmt.Sprintf("gateway:%s", reading.Gateway),
			)

			ddlambda.MetricWithTimestamp(
				"sensorpush.sensor.humidity",
				reading.Humidity,
				reading.Timestamp,
				fmt.Sprintf("name:%s", sensor.Name),
				fmt.Sprintf("type:%s", sensor.Type),
				fmt.Sprintf("gateway:%s", reading.Gateway),
			)
		}
	}

	return nil
}
