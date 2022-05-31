package validacion

import (
	"encoding/json"
	"io/ioutil"
)

func LeerArchivo(archivo string) ([]byte, error) {
	data, err := ioutil.ReadFile(archivo)
	return data, err
}

func DeserializarEntradaValidacion(data []byte) (*Entrada_validacion, error) {
	var h *Entrada_validacion
	err := json.Unmarshal(data, &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func SerializarSalidaValidacion(h *Salida_validacion) ([]byte, error) {
	data, err := json.MarshalIndent(h, "", " ")
	if err != nil {
		return nil, err
	}
	return data, err
}

