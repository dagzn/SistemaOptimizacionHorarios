package main

import (
	"fmt"
	obj "proyecto-horarios/objetos"
	"proyecto-horarios/utils"
	"proyecto-horarios/solucion"
	"proyecto-horarios/validacion"
)

func probarSolucion(archivo string) (*obj.Salida_horario){
	data, err := utils.LeerArchivo(archivo)
	if err != nil {
		return &obj.Salida_horario{
			Error: err.Error(),
		}
	}

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

func probarValidacion(archivo string) (*obj.Salida_validacion){
	data, err := utils.LeerArchivo(archivo)
	if err != nil {
		return &obj.Salida_validacion{
			Error: err.Error(),
		}
	}

	entradaValidacion, err := utils.DeserializarEntradaValidacion(data)
	if err != nil {
		return &obj.Salida_validacion{
			Error: err.Error(),
		}
	}

	errores, err := validacion.ValidarFormatoEntradaValidacion(entradaValidacion)
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
			salida := probarSolucion("archivos_pruebas/"+archivo)
			content, err := utils.SerializarSalidaHorario(salida)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(content))
		} else if opc == 2 {
			salida := probarValidacion("archivos_pruebas/"+archivo)
			content, err := utils.SerializarSalidaValidacion(salida)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(content))
		}
	}
}
