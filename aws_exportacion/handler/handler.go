package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"proyecto-horarios/exportacion"
	"proyecto-horarios/utils"

	"github.com/aws/aws-lambda-go/events"
)

func obtenerHeaders() map[string]string {
	return map[string]string{
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
		respuesta.Headers["Content-Type"] = "application/json"

		respuesta.Body = fmt.Sprintf(`
		{
			"distribuciones": null,
			"error": "%s",
			"logs": null
		}
		`, err.Error())

		respuesta.StatusCode = http.StatusOK
		return respuesta, nil
	}

	respuesta.Headers["Content-Type"] = "application/pdf"

	respuesta.IsBase64Encoded = true
	respuesta.Body = string(cadenaCodificada)
	respuesta.StatusCode = http.StatusOK

	return respuesta, nil
}
