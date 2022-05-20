package main

import (
	"fmt"
	"strconv"
	"time"
	"math/rand"
	"math"
	"os"
	"encoding/json"
	obj "proyecto-horarios/objetos"
	"proyecto-horarios/utils"
	"proyecto-horarios/solucion"
	"proyecto-horarios/validacion"
	"proyecto-horarios/exportacion"
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

func probarExportacion(archivo string) (string) {
	data, err := utils.LeerArchivo(archivo)
	if err != nil {
		panic(err)
	}

	entradaExportacion, err := utils.DeserializarEntradaExportacion(data)
	if err != nil {
		panic(err)
	}

	cadena, err := exportacion.ExportarHorario(entradaExportacion, "./")
	if err != nil {
		panic(err)
	}

	return cadena
}

func TestArchivo(){
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
		} else if opc == 3 {
			cadena := probarExportacion("archivos_pruebas/"+archivo)
			fmt.Print(cadena)
		}
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getLimInf(numClases, todas int) int {
	if todas + 1 < numClases {
		return numClases - todas - 1
	}
	return 1
}

func main(){

	var cantProfesores, cantMaterias, cantBloques, numClases, numAristas, idProceso int
	fmt.Printf("ID:\n")
	fmt.Scanf("%d", &idProceso)
	fmt.Printf("Profesores:\n")
	fmt.Scanf("%d", &cantProfesores)
	fmt.Printf("Materias:\n")
	fmt.Scanf("%d", &cantMaterias)
	fmt.Printf("Bloques:\n")
	fmt.Scanf("%d", &cantBloques)
	fmt.Printf("Numero de clases:\n")
	fmt.Scanf("%d", &numClases)
	fmt.Printf("Num aristas por profe (cada lado):\n")
	fmt.Scanf("%d", &numAristas)

	if cantProfesores > numClases {
		panic(fmt.Errorf("Numero de clases debe ser mayor o igual a num de profesores"))
	}

	if numClases > cantProfesores * numAristas {
		panic(fmt.Errorf("Numero de clases debe ser menor o igual a cantidad de conexiones"))
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	start := time.Now()

	var profesores []obj.Profesor
	totalClases := numClases
	for i := 0; i < cantProfesores; i++ {
		quedan := cantProfesores - (i+1)
		disponibles := numClases - quedan
		limiteSuperior := minInt(disponibles, numAristas)
		todasRestantes := quedan * numAristas
		limiteInferior := getLimInf(numClases, todasRestantes)
		clases := r1.Intn(limiteSuperior - limiteInferior + 1) + limiteInferior

		numClases = numClases - clases
		id_profe := strconv.Itoa(i+1)

		var prefMaterias []obj.Pref_materia
		m_usadas := make(map[int]bool)
		for j := 0; j < numAristas; j++ {
			var id_m int
			for true {
				id_m = r1.Intn(cantMaterias) + 1
				if _, ok := m_usadas[id_m]; !ok {
					m_usadas[id_m] = true
					break
				}
			}

			shiftSz := r1.Intn(10)
			pref := 1 << shiftSz
			materia := obj.Pref_materia{
				Id: strconv.Itoa(id_m),
				Limite: 5,
				Preferencia: pref,
			}

			prefMaterias = append(prefMaterias, materia)
		}

		var prefBloques []obj.Pref_bloque
		b_usados := make(map[int]bool)
		for j := 0; j < numAristas; j++ {
			var id_b int
			for true {
				id_b = r1.Intn(cantBloques) + 1
				if _, ok := b_usados[id_b]; !ok {
					b_usados[id_b] = true
					break
				}
			}

			shiftSz := r1.Intn(10)
			pref := 1 << shiftSz
			bloque := obj.Pref_bloque{
				Id: strconv.Itoa(id_b),
				Preferencia: pref,
			}

			prefBloques = append(prefBloques, bloque)
		}

		profe := obj.Profesor{
			Id: id_profe,
			Nombre: id_profe,
			Clases: clases,
			Materias: prefMaterias,
			Bloques: prefBloques,
		}

		profesores = append(profesores, profe)
	}

	var materias []obj.Materia
	for i := 0; i < cantMaterias; i++ {
		tocan := totalClases / cantMaterias
		if modu := int(math.Mod(float64(totalClases), float64(cantMaterias))); i < modu {
			tocan = tocan + 1
		}

		id_materia := strconv.Itoa(i + 1)
		mm := obj.Materia {
			Id: id_materia,
			Nombre: id_materia,
			Cantidad: tocan,
		}

		materias = append(materias, mm)
	}

	modulo := obj.Modulo{
		Dia: "Lunes",
		Entrada: "07:00",
		Salida: "10:00",
	}

	var modulos []obj.Modulo
	for i:=0; i < 5; i++ {
		modulos = append(modulos, modulo)
	}

	var bloques []obj.Bloque
	for i:=0; i < cantBloques; i++ {
		id_bloque := strconv.Itoa(i+1)

		bloque := obj.Bloque{
			Id: id_bloque,
			Nombre: id_bloque,
			Modulos: modulos,
		}

		bloques = append(bloques, bloque)
	}

	entradaHorario := &obj.Entrada_horario{
		Salones: 200000,
		Materias: materias,
		Bloques: bloques,
		Profesores: profesores,
	}

	data, err := json.MarshalIndent(entradaHorario, "", " ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("peticiones/peticion"+strconv.Itoa(idProceso)+".json", data, 0644)
	if err != nil {
		panic(err)
	}

	errores, err := validacion.ValidarFormatoEntradaHorario(entradaHorario)
	if err != nil {
		for _, x := range errores {
			fmt.Println(x)
		}
		return
	}

	salida, err := solucion.GenerarHorario(entradaHorario)
	if err != nil {
		panic(err)
	}

	fmt.Println("Terminamos el proceso!")

	bytes, err := utils.SerializarSalidaHorario(salida)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("resultados/resultado"+strconv.Itoa(idProceso)+".json", bytes, 0644)
	if err != nil {
		panic(err)
	}

	duration := time.Since(start)
	fmt.Println("Profesores: ", len(profesores))
	fmt.Println("Materias: ", len(materias))
	fmt.Println("Bloques: ", len(bloques))
	fmt.Println("Clases totales: ", totalClases)
	nodos := 2 + len(materias) + 2* len(profesores) + len(bloques)
	fmt.Println("Nodos creados: ", nodos)
	aristas := len(materias) + len(bloques) + len(profesores) + 2 *len(profesores)* numAristas
	fmt.Println("Aristas creadas: ", aristas)
	fmt.Println("Tiempo de ejecucion: ", duration)
}
