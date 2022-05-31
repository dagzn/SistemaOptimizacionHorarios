package solucion

import (
	"encoding/json"
	"io/ioutil"
)

func LeerArchivo(archivo string) ([]byte, error) {
	data, err := ioutil.ReadFile(archivo)
	return data, err
}

func DeserializarEntradaHorario(data []byte) (*Entrada_horario, error) {
	var h *Entrada_horario
	err := json.Unmarshal(data, &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func SerializarSalidaHorario(h *Salida_horario) ([]byte, error) {
	data, err := json.MarshalIndent(h, "", " ")
	if err != nil {
		return nil, err
	}
	return data, err
}

