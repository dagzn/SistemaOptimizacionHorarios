package main

import (
	"proyecto-horarios/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.AtenderPeticion)
}
