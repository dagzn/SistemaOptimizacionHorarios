package main

import (
	"fmt"
	sol "proyecto-horarios/solucion"
	val "proyecto-horarios/validacion"
	exp "proyecto-horarios/exportacion"
)

func probarSolucion(archivo string) (*sol.Salida_horario){
	data, err := sol.LeerArchivo(archivo)
	if err != nil {
		return &sol.Salida_horario{
			Error: err.Error(),
		}
	}

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

func probarValidacion(archivo string) (*val.Salida_validacion){
	data, err := val.LeerArchivo(archivo)
	if err != nil {
		return &val.Salida_validacion{
			Error: err.Error(),
		}
	}

	entradaValidacion, err := val.DeserializarEntradaValidacion(data)
	if err != nil {
		return &val.Salida_validacion{
			Error: err.Error(),
		}
	}

	errores, err := val.ValidarFormatoEntradaValidacion(entradaValidacion)
	if err != nil {
		return &val.Salida_validacion{
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

	salida := val.ValidarHorario(entradaValidacion)

	return salida
}

func probarExportacion(archivo string) (string) {
	data, err := exp.LeerArchivo(archivo)
	if err != nil {
		panic(err)
	}

	entradaExportacion, err := exp.DeserializarEntradaExportacion(data)
	if err != nil {
		panic(err)
	}

	cadena, err := exp.ExportarHorario(entradaExportacion, "./")
	if err != nil {
		panic(err)
	}

	return cadena
}

func main(){
	var archivo string
	opc := 0
	for true {
		fmt.Printf("\n1.- Probar servicio de solucion.\n2.- Probar servicio de validacion.\n3.- Servicio de exportacion\n4.- Salir\n");
		fmt.Scanf("%d", &opc)
		if opc == 4 {
			break
		}
		fmt.Printf("Nombre del archivo:\n")
		fmt.Scanf("%s", &archivo)
		if opc == 1 {
			salida := probarSolucion("archivos_pruebas/"+archivo)
			content, err := sol.SerializarSalidaHorario(salida)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(content))
		} else if opc == 2 {
			salida := probarValidacion("archivos_pruebas/"+archivo)
			content, err := val.SerializarSalidaValidacion(salida)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(content))
		} else if opc == 3 {
			cadena := probarExportacion("archivos_pruebas/"+archivo)
			fmt.Print(cadena)
		}
	}
}
