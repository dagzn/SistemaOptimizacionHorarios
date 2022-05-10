package peticion

import (
	"encoding/base64"
	"fmt"
	"net/http"
	obj "proyecto-horarios/objetos"
	"proyecto-horarios/solucion"
	"proyecto-horarios/utils"
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

func probarSolucion(data []byte) *obj.Salida_horario {
	entradaHorario, err := utils.DeserializarEntradaHorario(data)
	if err != nil {
		return &obj.Salida_horario{
			Error: err.Error(),
		}
	}

	errores, err := validacion.ValidarFormatoEntradaHorario(entradaHorario)
	if err != nil {
		return &obj.Salida_horario{
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

	salida, err := solucion.GenerarHorario(entradaHorario)
	if err != nil {
		return &obj.Salida_horario{
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
	content, err := utils.SerializarSalidaHorario(salida)
	if err != nil {
		respuesta.Body = fmt.Sprintf(`
			"Distribuciones": null,
			"Error": "%s",
			"Logs": null
		`, err.Error())

		respuesta.StatusCode = http.StatusOK
		return respuesta, nil
	}

	respuesta.Body = string(content)
	respuesta.StatusCode = http.StatusOK

	return respuesta, nil
}
