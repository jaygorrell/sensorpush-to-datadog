# Sensorpush to Datadog

## Requirements
- [Serverless Framework](https://www.serverless.com/framework/docs)
- A [Datadog](https://www.datadoghq.com/) Account and API Key
- An AWS Account with Admin access configured in a local profile
- A SensorPush account (and ideally some sensors!)

## Configure

Start by making a copy of `serverless.example.yml` to `serverless.yml`. In your new file, set your `awsProfile` 
(if not using `default`) and your `awsRegion`.

### Datadog Secret
> For other options, see the [serverless-plugin-datadog](https://www.serverless.com/plugins/serverless-plugin-datadog) page.

1. Using AWS Secrets Manager, create a plaintext secret containing the API key to your Datadog environment.
2. Set the `datadogSecret` value to the ARN of the secret created above by replacing the `DATADOG_SECRET_ARN` placeholder.

### Sensorpush Secret

1. Using AWS Secrets Manager, create a key/value secret containing the attributes `email` and `password` with your SensorPush credentials.
2. Set the `sensorpushSecret` value to the ARN of the secret created above by replacing the `SENSORPUSH_SECRET_ARN` placeholder.

## Build

Mac:
```shell
env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./bin/bootstrap -tags lambda.norpc handlers/temperature.go
```

Win (PowerShell): 
```shell
$Env:GOOS="linux"; $Env:GOARCH="amd64"; go build -ldflags="-s -w" -o ./bootstrap -tags lambda.norpc cmd/temperature/main.go
```

## Deploy
> Any other stage name can be used (`dev` is used by default) to create a test environment with the cron schedule disabled. 

`serverless deploy --stage prod`
