package peticion

import (
	"encoding/base64"
	"fmt"
	"net/http"
	sol "proyecto-horarios/solucion"

	"github.com/aws/aws-lambda-go/events"
)

func obtenerHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, HEAD, OPTIONS, POST",
	}
}

func probarSolucion(data []byte) *sol.Salida_horario {
	entradaHorario, err := sol.DeserializarEntradaHorario(data)
	if err != nil {
		return &sol.Salida_horario{
			Error: err.Error(),
		}
	}

	errores, err := sol.ValidarFormatoEntradaHorario(entradaHorario)
	if err != nil {
		return &sol.Salida_horario{
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

	salida, err := sol.GenerarHorario(entradaHorario)
	if err != nil {
		return &sol.Salida_horario{
			Error: err.Error(),
		}
	}

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

	salida := probarSolucion([]byte(body))
	content, err := sol.SerializarSalidaHorario(salida)
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
