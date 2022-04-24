package validacion

import (
	"fmt"
	obj "proyecto-horarios/objetos"
)

const (
	errorValidacion = "Se encontraron errores al momento de validar el horario. Cheque los logs para mas informacion."
	errorClasesMateria = "La materia '%s' tiene que impartirse %d veces y tiene %d clases asignadas."
	errorClasesProfesor = "El profesor %s requiere de %d clases en total y se le asignaron %d."
	errorInterseccionBloques = "El profesor %s tiene mas de una clase asignada para el bloque %s."
	errorLimiteBloque = "El bloque %s tiene un maximo de %d salones y se le asignaron %d clases."
	errorLimiteMateria = "El limite impuesto para el profesor %s para dar la materia %s es de %d clases pero se le asignaron %d clases."
	errorMateriaNoAsignada = "La materia %s no puede ser impartida por el profesor %s."
	errorBloqueNoAsignado = "El profesor %s no puede impartir clases en el bloque %s."
)

var (
	materiasAsignadas map[int][]int
	bloquesAsignados map[int][]int
	nombreMateria map[int]string
	nombreBloque map[int]string
)

/*
	Validar que cada materia sea dada una cierta cantidad de veces.
*/
func validarMaterias(materias []obj.Materia, distribuciones []obj.Distribucion, validaciones map[string]int) ([]error) {
	if _, ok := validaciones["ClasesMateria"]; !ok {
		return nil
	}

	clasesAsignadas := make(map[string]int)
	for _, d := range distribuciones {
		for _, a := range d.Asignaciones {
			materia := a.Materia
			actual := clasesAsignadas[materia]
			clasesAsignadas[materia] = actual + 1
		}
	}

	var errores []error
	for _, m := range materias {
		asignadas := clasesAsignadas[m.Nombre]
		if m.Cantidad != asignadas {
			errores = append(errores, fmt.Errorf(errorClasesMateria, m.Nombre, m.Cantidad, asignadas))
		}
	}

	return errores
}

/*
	Validar que cada bloque cumpla con con limite asignado de salones.
*/
func validarBloques(salonesDisponibles int, distribuciones []obj.Distribucion, validaciones map[string]int) ([]error) {
	if _, ok := validaciones["LimiteBloque"]; !ok {
		return nil;
	}

	var errores []error
	for _, d := range distribuciones {
		salonesOcupados := len(d.Asignaciones)
		if salonesDisponibles < salonesOcupados {
			errores = append(errores, fmt.Errorf(errorLimiteBloque, d.Bloque.Nombre, salonesDisponibles, salonesOcupados))
		}
	}

	return errores
}

func unique(intSlice []int) []int {
    keys := make(map[int]bool)
    list := []int{}
    for _, entry := range intSlice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
}

// Validar que un profesor no tenga mas de una clase en el mismo bloque.
func validarInterseccionBloques(distribuciones []obj.Distribucion, validaciones map[string]int) ([]error) {
	if _, ok := validaciones["InterseccionBloques"]; !ok {
		return nil
	}

	var errores []error
	for _, d := range distribuciones {
		duplicados := make(map[int]bool)
		for _, a := range d.Asignaciones {
			profesor := a.Id_profesor
			if ok := duplicados[profesor]; ok {
				errores = append(errores, fmt.Errorf(errorInterseccionBloques, a.Profesor, d.Bloque.Nombre))
			} else {
				duplicados[profesor] = true
			}
		}
	}
	return errores
}


// Validar que se cumplan con las clases requeridas para cada profesor.
func validarClasesProfesor(profesores []obj.Profesor, validaciones map[string]int) ([]error) {
	if _, ok := validaciones["ClasesProfesor"]; !ok {
		return nil
	}

	var errores []error
	for _, p := range profesores {
		asignadas := len(materiasAsignadas[p.Id])
		if asignadas != p.Clases {
			errores = append(errores, fmt.Errorf(errorClasesProfesor, p.Nombre, p.Clases, asignadas))
		}
	}
	return errores
}

// Validar que el profesor no de mas veces una materia que las deseadas
func validarLimiteMateria(profesores []obj.Profesor, validaciones map[string]int) ([]error) {
	if _, ok := validaciones["LimiteMateria"]; !ok {
		return nil
	}

	var errores []error
	for _, p := range profesores {
		vecesAsignada := make(map[int]int)
		for _, id := range materiasAsignadas[p.Id] {
			vecesAsignada[id] = vecesAsignada[id] + 1
		}

		materiasAsignadas[p.Id] = unique(materiasAsignadas[p.Id])
		for _, m := range p.Materias {
			if veces, ok := vecesAsignada[m.Id]; ok {
				if veces > m.Limite {
					errores = append(errores, fmt.Errorf(errorLimiteMateria, p.Nombre, nombreMateria[m.Id], m.Limite, veces))
				}
			}
		}
	}
	return errores
}

// Validar de que un profesor no tenga asignada materia que no puede dar.
func validarMateriaNoAsignada(profesores []obj.Profesor, validaciones map[string]int) ([]error) {
	if _, ok := validaciones["MateriaNoAsignada"]; !ok {
		return nil
	}

	var errores []error
	for _, p := range profesores {
		preferencias := make(map[int]int)
		for _, m := range p.Materias {
			preferencias[m.Id] = 1
		}

		for _, id := range materiasAsignadas[p.Id] {
			if _, ok := preferencias[id]; !ok {
				errores = append(errores, fmt.Errorf(errorMateriaNoAsignada, nombreMateria[id], p.Nombre))
			}
		}
	}
	return errores
}

// Validar que el profesor solo tenga asignados bloques en los que si pueda trabajar.
func validarBloqueNoAsignado(profesores []obj.Profesor, validaciones map[string]int) ([]error) {
	if _, ok := validaciones["BloqueNoAsignado"]; !ok {
		return nil
	}

	var errores []error
	for _, p := range profesores {
		preferencias := make(map[int]int)
		for _, b := range p.Bloques {
			preferencias[b.Id] = 1
		}

		bloquesAsignados[p.Id] = unique(bloquesAsignados[p.Id])

		for _, id := range bloquesAsignados[p.Id] {
			if _, ok := preferencias[id]; !ok {
				errores = append(errores, fmt.Errorf(errorBloqueNoAsignado, p.Nombre, nombreBloque[id]))
			}
		}
	}
	return errores
}

func validarProfesores(profesores []obj.Profesor, distribuciones []obj.Distribucion, validaciones map[string]int) ([]error) {
	var errores []error

	err := validarInterseccionBloques(distribuciones, validaciones)
	if len(err) > 0 {
		errores = append(errores, err...)
	}

	err = validarClasesProfesor(profesores, validaciones)
	if len(err) > 0 {
		errores = append(errores, err...)
	}

	err = validarLimiteMateria(profesores, validaciones)
	if len(err) > 0 {
		errores = append(errores, err...)
	}

	err = validarMateriaNoAsignada(profesores, validaciones)
	if len(err) > 0 {
		errores = append(errores, err...)
	}

	err = validarBloqueNoAsignado(profesores, validaciones)
	if len(err) > 0 {
		errores = append(errores, err...)
	}

	return errores
}

func llenarInformacionProfesores(distribuciones []obj.Distribucion) {
	materiasAsignadas = make(map[int][]int)
	bloquesAsignados = make(map[int][]int)
	nombreMateria = make(map[int]string)
	nombreBloque = make(map[int]string)

	for _, d := range distribuciones {
		for _, a := range d.Asignaciones {
			profesor := a.Id_profesor
			materiasAsignadas[profesor] = append(materiasAsignadas[profesor], a.Id_materia)
			bloquesAsignados[profesor] = append(bloquesAsignados[profesor], d.Bloque.Id)
			nombreMateria[a.Id_materia] = a.Materia
		}
		nombreBloque[d.Bloque.Id] = d.Bloque.Nombre
	}
}

func validar(horario *obj.Entrada_validacion) ([]error) {
	var errores []error

	listaValidaciones := horario.Validaciones
	validaciones := make(map[string]int)
	for _, v := range listaValidaciones {
		validaciones[v] = 1
	}

	erroresMaterias := validarMaterias(horario.Materias, horario.Distribuciones, validaciones)
	if len(erroresMaterias) > 0 {
		errores = append(errores, erroresMaterias...)
	}

	erroresBloques := validarBloques(horario.Salones, horario.Distribuciones, validaciones)
	if len(erroresBloques) > 0 {
		errores = append(errores, erroresBloques...)
	}

	llenarInformacionProfesores(horario.Distribuciones)

	erroresProfesores := validarProfesores(horario.Profesores, horario.Distribuciones, validaciones)
	if len(erroresProfesores) > 0 {
		errores = append(errores, erroresProfesores...)
	}

	return errores
}

func ValidarHorario(horario *obj.Entrada_validacion) (*obj.Salida_validacion){
	errores := validar(horario)
	if len(errores) > 1 {
		err := fmt.Errorf(errorValidacion);
		salida := &obj.Salida_validacion{
			Error: err.Error(),
			Logs: func(errores []error) []string {
				var ret []string
				for _, e := range errores {
					ret = append(ret, e.Error())
				}
				return ret
			}(errores),
		}

		return salida
	}

	return &obj.Salida_validacion{}
}
