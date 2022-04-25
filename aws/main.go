package main

import (
	"proyecto-horarios/peticion"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(peticion.AtenderPeticion)
}
