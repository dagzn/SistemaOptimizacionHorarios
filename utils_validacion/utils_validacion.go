package utils_validacion

import (
	"encoding/json"
	"io/ioutil"
	obj "proyecto-horarios/objetos_validacion"
)

func LeerArchivo(archivo string) ([]byte, error) {
	data, err := ioutil.ReadFile(archivo)
	return data, err
}

func DeserializarEntradaValidacion(data []byte) (*obj.Entrada_validacion, error) {
	var h *obj.Entrada_validacion
	err := json.Unmarshal(data, &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func SerializarSalidaValidacion(h *obj.Salida_validacion) ([]byte, error) {
	data, err := json.MarshalIndent(h, "", " ")
	if err != nil {
		return nil, err
	}
	return data, err
}

