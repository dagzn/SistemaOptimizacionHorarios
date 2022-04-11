package utils

import (
	"encoding/json"
	"io/ioutil"
	obj "proyecto-horarios/objetos"
)

func LeerArchivo(archivo string) ([]byte, error) {
	data, err := ioutil.ReadFile(archivo)
	return data, err
}

func DeserializarEntradaHorario(data []byte) (*obj.Entrada_horario, error){
	var h *obj.Entrada_horario
	err := json.Unmarshal(data, &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func SerializarSalidaHorario(h *obj.Salida_horario) ([]byte, error) {
	data, err := json.MarshalIndent(h, "", " ")
	if err != nil {
		return nil, err
	}
	return data, err
}
