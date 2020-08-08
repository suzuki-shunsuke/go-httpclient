# go-httpclient

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/suzuki-shunsuke/go-httpclient/httpclient)
[![Build Status](https://cloud.drone.io/api/badges/suzuki-shunsuke/go-httpclient/status.svg)](https://cloud.drone.io/suzuki-shunsuke/go-httpclient)
[![Test Coverage](https://api.codeclimate.com/v1/badges/3fa33b6aa6830f36406e/test_coverage)](https://codeclimate.com/github/suzuki-shunsuke/go-httpclient/test_coverage)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/go-httpclient)](https://goreportcard.com/report/github.com/suzuki-shunsuke/go-httpclient)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/go-httpclient.svg)](https://github.com/suzuki-shunsuke/go-httpclient)
[![GitHub tag](https://img.shields.io/github/tag/suzuki-shunsuke/go-httpclient.svg)](https://github.com/suzuki-shunsuke/go-httpclient/releases)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/go-httpclient/master/LICENSE)

Go's simple HTTP client.

## Assumption

* The format of request and response body is JSON
* If the response status code is greater equal than 300, it is treated as an error response

## License

[MIT](LICENSE)
