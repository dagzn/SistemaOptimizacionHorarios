package main

import (
	"fmt"
	"proyecto-horarios/utils"
	"proyecto-horarios/solucion"
)

func probarSolucion(){
	data, err := utils.LeerArchivo("entrada_grafo.json")
	if err != nil {
		panic(err)
	}

	entradaHorario, err := utils.DeserializarEntradaHorario(data)
	if err != nil {
		panic(err)
	}

	salida, err := solucion.GenerarHorario(entradaHorario)
	if err != nil {
		panic(err)
	}

	bytes, err := utils.SerializarSalidaHorario(salida)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}

func main(){
	probarSolucion()
}
