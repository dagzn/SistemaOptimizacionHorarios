package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"

	val "proyecto-horarios/validacion"

	"github.com/aws/aws-lambda-go/events"
)

func obtenerHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, HEAD, OPTIONS, POST",
	}
}

func probarValidacion(data []byte) *val.Salida_validacion {

	entradaValidacion, err := val.DeserializarEntradaValidacion(data)
	if err != nil {
		return &val.Salida_validacion{
			Error: err.Error(),
		}
	}

	errores, err := val.ValidarFormatoEntradaValidacion(entradaValidacion)
	if err != nil {
		return &val.Salida_validacion{
			Error: err.Error(),
			Logs: func(errores []error) []string {
				var ret []string
				for _, e := range errores {
					ret = append(ret, e.Error())
				}
				return ret
			}(errores),
		}
	}

	salida := val.ValidarHorario(entradaValidacion)

	return salida
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

	salida := probarValidacion([]byte(body))

	content, err := val.SerializarSalidaValidacion(salida)
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
