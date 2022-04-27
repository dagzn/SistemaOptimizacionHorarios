package handler

import (
	"encoding/base64"
	"net/http"
	"proyecto-horarios/exportacion"
	"proyecto-horarios/utils"

	"github.com/aws/aws-lambda-go/events"
)

func obtenerHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                 "application/pdf",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, HEAD, OPTIONS, POST",
	}
}

func probarExportacion(data []byte) (string, error) {
	entradaExportacion, err := utils.DeserializarEntradaExportacion(data)
	if err != nil {
		return "", err
	}

	cadenaCodificada, err := exportacion.ExportarHorario(entradaExportacion, "/tmp/")
	if err != nil {
		return "", err
	}

	return cadenaCodificada, nil
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

	cadenaCodificada, err := probarExportacion([]byte(body))

	if err != nil {
		respuesta.StatusCode = http.StatusInternalServerError
		return respuesta, err
	}

	respuesta.IsBase64Encoded = true
	respuesta.Body = string(cadenaCodificada)
	respuesta.StatusCode = http.StatusOK

	return respuesta, nil
}
