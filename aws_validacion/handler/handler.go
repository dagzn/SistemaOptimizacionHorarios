package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"

	obj "proyecto-horarios/objetos_validacion"
	utils "proyecto-horarios/utils_validacion"
	formato "proyecto-horarios/formato_validacion"
	"proyecto-horarios/validacion"

	"github.com/aws/aws-lambda-go/events"
)

func obtenerHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, HEAD, OPTIONS, POST",
	}
}

func probarValidacion(data []byte) *obj.Salida_validacion {

	entradaValidacion, err := utils.DeserializarEntradaValidacion(data)
	if err != nil {
		return &obj.Salida_validacion{
			Error: err.Error(),
		}
	}

	errores, err := formato.ValidarFormatoEntradaValidacion(entradaValidacion)
	if err != nil {
		return &obj.Salida_validacion{
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

	salida := validacion.ValidarHorario(entradaValidacion)

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

	content, err := utils.SerializarSalidaValidacion(salida)
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
