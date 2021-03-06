package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"

	exp "proyecto-horarios/exportacion"

	"github.com/aws/aws-lambda-go/events"
)

func obtenerHeaders() map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, HEAD, OPTIONS, POST",
	}
}

func probarExportacion(data []byte) (string, string, *exp.Salida_exportacion_fallida, error) {
	entradaExportacion, err := exp.DeserializarEntradaExportacion(data)
	if err != nil {
		return "", "", nil, err
	}

	errores, err := exp.ValidarFormatoEntradaExportacion(entradaExportacion)
	if err != nil {
		return "", "", &exp.Salida_exportacion_fallida{
			Error: err.Error(),
			Logs: func(errores []error) []string {
				var ret []string
				for _, e := range errores {
					ret = append(ret, e.Error())
				}
				return ret
			}(errores),
		}, nil
	}

	cadenaCodificada, err := exp.ExportarHorario(entradaExportacion, "/tmp/")
	if err != nil {
		return "", "", nil, err
	}

	if entradaExportacion.Tipo == "Individual" {
		return cadenaCodificada, "application/zip", nil, nil
	}

	return cadenaCodificada, "application/pdf", nil, nil
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

	cadenaCodificada, contentType, salidaExportacionFallida, err := probarExportacion([]byte(body))

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

	if salidaExportacionFallida != nil {
		respuesta.Headers["Content-Type"] = "application/json"

		content, err := exp.SerializarSalidaExportacionFallida(salidaExportacionFallida)

		if err != nil {
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

		respuesta.Body = string(content)
		respuesta.StatusCode = http.StatusOK

		return respuesta, nil
	}

	respuesta.Headers["Content-Type"] = contentType

	respuesta.IsBase64Encoded = true
	respuesta.Body = string(cadenaCodificada)
	respuesta.StatusCode = http.StatusOK

	return respuesta, nil
}
