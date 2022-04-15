package validacion

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	obj "proyecto-horarios/objetos"
)

const(
	errorGeneral = "Existe un error en el formato de la peticion. Cheque los logs para obtener mas informacion.\n"
	errorRequired = "El campo %s es obligatorio.\nCampo: %s\n"
	errorMin = "El arreglo %s no cuenta con la longitud esperada.\n%s"
	errorGte = "El valor del campo %s es menor al valor esperado.\n%s"
	errorOneof = "El valor del campo %s no se encuentra entre la lista de valores esperados.\n%s"
	infoCampo = "Campo: %s\nValor esperado: %+v\nValor real: %+v\n"
)


func ValidarFormatoEntradaHorario(h *obj.Entrada_horario) ([]error, error) {
	validate := validator.New()
	err := validate.Struct(h)

	if err != nil {
		var errores []error
		for _, err := range err.(validator.ValidationErrors) {
			tag := err.ActualTag()
			var errFmt error
			switch tag {
				case "required":
					errFmt = fmt.Errorf(errorRequired,err.StructField(), err.StructNamespace())
				case "gte":
					errFmt = fmt.Errorf(errorGte, err.StructField(),fmt.Sprintf(infoCampo, err.StructNamespace(), err.Param(), err.Value()))
				case "oneof":
					errFmt = fmt.Errorf(errorOneof, err.StructField(),fmt.Sprintf(infoCampo, err.StructNamespace(), err.Param(), err.Value()))
				case "min":
					errFmt = fmt.Errorf(errorMin, err.StructField(),fmt.Sprintf(infoCampo, err.StructNamespace(), err.Param(), err.Value()))
			}

			errores = append(errores, errFmt)
		}
		return errores, fmt.Errorf(errorGeneral)
	}

	return nil, nil
}

func ValidarFormatoEntradaValidacion(h *obj.Entrada_validacion) ([]error, error) {
	validate := validator.New()
	err := validate.Struct(h)

	if err != nil {
		var errores []error
		for _, err := range err.(validator.ValidationErrors) {
			tag := err.ActualTag()
			var errFmt error
			switch tag {
				case "required":
					errFmt = fmt.Errorf(errorRequired,err.StructField(), err.StructNamespace())
				case "gte":
					errFmt = fmt.Errorf(errorGte, err.StructField(),fmt.Sprintf(infoCampo, err.StructNamespace(), err.Param(), err.Value()))
				case "oneof":
					errFmt = fmt.Errorf(errorOneof, err.StructField(),fmt.Sprintf(infoCampo, err.StructNamespace(), err.Param(), err.Value()))
				case "min":
					errFmt = fmt.Errorf(errorMin, err.StructField(),fmt.Sprintf(infoCampo, err.StructNamespace(), err.Param(), err.Value()))
			}

			errores = append(errores, errFmt)
		}
		return errores, fmt.Errorf(errorGeneral)
	}

	return nil, nil
}
