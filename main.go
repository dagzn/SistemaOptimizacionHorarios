package main

import (
	"fmt"
	objs "proyecto-horarios/objetos_solucion"
	objv "proyecto-horarios/objetos_validacion"
	utilss "proyecto-horarios/utils_solucion"
	utilsv "proyecto-horarios/utils_validacion"
	utilse "proyecto-horarios/utils_exportacion"
	fmts "proyecto-horarios/formato_solucion"
	fmtv "proyecto-horarios/formato_validacion"
	"proyecto-horarios/solucion"
	"proyecto-horarios/validacion"
	"proyecto-horarios/exportacion"
)

func probarSolucion(archivo string) (*objs.Salida_horario){
	data, err := utilss.LeerArchivo(archivo)
	if err != nil {
		return &objs.Salida_horario{
			Error: err.Error(),
		}
	}

	entradaHorario, err := utilss.DeserializarEntradaHorario(data)
	if err != nil {
		return &objs.Salida_horario{
			Error: err.Error(),
		}
	}

	errores, err := fmts.ValidarFormatoEntradaHorario(entradaHorario)
	if err != nil {
		return &objs.Salida_horario{
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
		return &objs.Salida_horario{
			Error: err.Error(),
		}
	}

	return salida
}

func probarValidacion(archivo string) (*objv.Salida_validacion){
	data, err := utilsv.LeerArchivo(archivo)
	if err != nil {
		return &objv.Salida_validacion{
			Error: err.Error(),
		}
	}

	entradaValidacion, err := utilsv.DeserializarEntradaValidacion(data)
	if err != nil {
		return &objv.Salida_validacion{
			Error: err.Error(),
		}
	}

	errores, err := fmtv.ValidarFormatoEntradaValidacion(entradaValidacion)
	if err != nil {
		return &objv.Salida_validacion{
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

func probarExportacion(archivo string) (string) {
	data, err := utilse.LeerArchivo(archivo)
	if err != nil {
		panic(err)
	}

	entradaExportacion, err := utilse.DeserializarEntradaExportacion(data)
	if err != nil {
		panic(err)
	}

	cadena, err := exportacion.ExportarHorario(entradaExportacion, "./")
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
			content, err := utilss.SerializarSalidaHorario(salida)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(content))
		} else if opc == 2 {
			salida := probarValidacion("archivos_pruebas/"+archivo)
			content, err := utilsv.SerializarSalidaValidacion(salida)
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
