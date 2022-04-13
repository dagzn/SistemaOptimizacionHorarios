package main

import (
	"fmt"
	"proyecto-horarios/utils"
	"proyecto-horarios/solucion"
	"proyecto-horarios/validacion"
)

func probarSolucion(archivo string){
	data, err := utils.LeerArchivo(archivo)
	if err != nil {
		panic(err)
	}

	entradaHorario, err := utils.DeserializarEntradaHorario(data)
	if err != nil {
		panic(err)
	}

	errores := validacion.ValidarFormatoEntradaHorario(entradaHorario)
	if errores != nil {
		for _,e := range errores {
			fmt.Println(e.Error())
		}
		panic(fmt.Errorf("Existe un error en el formato de la peticion. Cheque los logs para obtener mas informacion.\n"))
	}

	salida, err := solucion.GenerarHorario(entradaHorario)
	if err != nil {
		panic(err)
	}

	content, err := utils.SerializarSalidaHorario(salida)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(content))
}

func probarValidacion(archivo string){
	_, err := utils.LeerArchivo(archivo)
	if err != nil {
		panic(err)
	}
}

func main(){
	var archivo string
	opc := 0
	for true {
		fmt.Printf("\n1.- Probar servicio de solucion.\n2.- Probar servicio de validacion.\n3.- Salir\n");
		fmt.Scanf("%d", &opc)
		if opc == 3 {
			break
		}
		fmt.Printf("Nombre del archivo:\n")
		fmt.Scanf("%s", &archivo)
		if opc == 1 {
			probarSolucion("archivos_pruebas/"+archivo)
		} else if opc == 2 {
			probarValidacion("archivos_pruebas/"+archivo)
		}
	}
}
