# About

This is an example client implementation written in GO to access the CyberVox platform API.

For more information, read the [API documentation](https://apidocs.cybervox.ai/).

# Getting Started

To run the example implementation:

```console
export CLIENT_ID=< provided client id >
export CLIENT_SECRET=< provided client secret >

# complete API
go run cmd/example/main.go

# text-to-speech only
go run cmd/tts/main.go "olá mundo"

# speech-to-text only
go run cmd/stt/main.go ola-mundo.wav
```

# Usage

[![GoDoc](https://godoc.org/github.com/cyberlabsai/cybervox-client-go?status.svg)](https://godoc.org/github.com/cyberlabsai/cybervox-client-go)

```console
go doc -all
```
