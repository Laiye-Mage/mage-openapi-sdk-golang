
## Introduction
This sample code demonstrate how to use Laiye mage API using golang. Full API doc is [here](https://mage.uibot.com.cn/docs/latest/docUnderstanding/backend/api.html).

## Get Started
Copy `env.go.sample` to `env.go` and put the public and secret key for each service. You can get these keys in mage's dashboard. See https://mage.uibot.com.cn/.

## How to use
Please see `api_test.go`. It should be quite straightforward, except you should define your response struct for each endpoint.


### Test all API and print results
```
go test -v

```


### Test a single API
```
go test -run TestNormalizeAddress
```
