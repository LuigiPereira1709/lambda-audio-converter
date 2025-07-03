package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"pitanguinha.com/audio-converter/handler"
)

func main() {
	lambda.Start(handler.Handler)
}
