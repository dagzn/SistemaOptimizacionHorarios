package peticion

import (
	"encoding/base64"
	"net/http"
	"proyecto-horarios/solucion"
	"proyecto-horarios/utils"

	"github.com/aws/aws-lambda-go/events"
)

func obtenerHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, HEAD, OPTIONS, POST",
	}
}

func AtenderPeticion(peticion events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	respuesta := events.APIGatewayProxyResponse{
		Headers: obtenerHeaders(),
	}

	var body string
	if peticion.IsBase64Encoded {
		bbody, _ := base64.StdEncoding.DecodeString(peticion.Body)
		body = string(bbody)
	} else {
		body = peticion.Body
	}

	entradaHorario, err := utils.DeserializarEntradaHorario([]byte(body))

	if err != nil {
		respuesta.StatusCode = http.StatusInternalServerError
		return respuesta, err
	}

	salida, err := solucion.GenerarHorario(entradaHorario)
	if err != nil {
		respuesta.StatusCode = http.StatusInternalServerError
		return respuesta, err
	}

	content, err := utils.SerializarSalidaHorario(salida)
	if err != nil {
		respuesta.StatusCode = http.StatusInternalServerError
		return respuesta, err
	}

	respuesta.Body = string(content)
	respuesta.StatusCode = http.StatusOK

	return respuesta, nil
}