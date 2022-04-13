package validacion

import (
	"fmt"
	obj "proyecto-horarios/objetos"
)

const (
	errorClasesMateria = "La materia '%s' tiene que impartirse %d veces y tiene %d clases asignadas."
	errorBloques = "El bloque %s tiene un maximo de %d salones y se le asignaron %d clases."
	errorClasesProfesor = "El profesor %s requiere de %d clases en total y se le asignaron %d."
	errorInterseccionBloques = "El profesor %s tiene mas de una clase asignada para el bloque %s."
	errorMateriaNoAsignada = "La materia %s no puede ser impartida por el profesor %s."
	errorLimiteMateria = "El limite impuesto para el profesor %s para dar la materia %s es de %d clases, pero se le asignaron %d clases."
	errorBloqueNoAsignado = "El profesor %s no puede impartir clases en el bloque %s."
)

/*
	Validar que cada materia sea dada una cierta cantidad de veces.
*/
func validarMaterias(materias []obj.Materia, distribuciones []obj.Distribucion) ([]error) {
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
func validarBloques(salonesDisponibles int, distribuciones []obj.Distribucion) ([]error) {
	var errores []error
	for _, d := range distribuciones {
		salonesOcupados := len(d.Asignaciones)
		if salonesDisponibles < salonesOcupados {
			errores = append(errores, fmt.Errorf(errorBloques, d.Bloque, salonesDisponibles, salonesOcupados))
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

func validarProfesores(profesores []obj.Profesor, distribuciones []obj.Distribucion) ([]error) {
	var errores []error
	/*
	materiasAsignadas := make(map[int][]int)
	bloquesAsignados := make(map[int][]int)
	nombreMateria := make(map[int]string)
	nombreBloque := make(map[int]string)

	for _, d := range distribuciones {
		duplicados := make(map[string]bool)
		for _, a := range d.Asignaciones {
			profesor := a.Profesor
			if duplicados[profesor] {
				// Validar que un profesor no tenga mas de una clase en el mismo bloque.
				errores = append(errores, fmt.Errorf(errorInterseccionBloques, profesor, d.Bloque))
			} else {
				duplicados[profesor] = true
			}

			materiasAsignadas[profesor] = append(materiasAsignadas[profesor], a.Id_materia)
			bloquesAsignados[profesor] = append(bloquesAsignados[profesor], d.Id_bloque)
			nombreMateria[a.Id_materia] = a.Materia
		}
		nombreBloque[d.Id_bloque] = d.Bloque
	}

	for _, p := range profesores {
		asignadas := len(materiasAsignadas[p.Nombre])
		if asignadas != p.Clases {
			// Validar que se cumplan con las clases requeridas para cada profesor.
			errores = append(errores, fmt.Errorf(errorClasesProfesor, p.Nombre, p.Clases, asignadas))
		}
	}

	for _, p := range profesores {
		vecesAsignada := make(map[int]int)
		preferencias := make(map[int]int)
		for _, m := range p.Materias {
			preferencias[m.Id] = 1
		}

		materiasAsignadas[p.Id] = unique(materiasAsignadas[p.Id])

		for _, id := range materiasAsignadas[p.Id] {
			if _, ok := preferencias[id]; ok {
				vecesAsignada[id] = vecesAsignada[id] + 1
			} else {
				// error de que tiene asignada materia que no
				errores = append(errores, fmt.Errorf(errorMateriaNoAsignada, nombreMateria[id], p.Nombre))
			}
		}

		for _, m := range p.Materias {
			if veces, ok := vecesAsignada[m.Id]; ok {
				if veces > m.Limite {
					// error que la da mas veces que las deseadas
					errores = append(errores, fmt.Errorf(errorLimiteMateria, p.Nombre, nombreMateria[m.Id], m.Limite, veces))
				}
			}
		}
	}

	// Que el profesor SOLO tenga asignados bloques en los que si pueda trabajar.
	for _, p := range profesores {
		preferencias := make(map[int]int)
		for _, b := range p.Bloques {
			preferencias[b.Id] = 1
		}

		for _, id := range bloquesAsignados[p.Id] {
			if _, ok := preferencias[id]; !ok {
				// el bloque asignado no esta en sus preferencias.
				errores = append(errores, fmt.Errorf(errorBloqueNoAsignado, p.Nombre, nombreBloque[id]))
			}
		}
	}
	*/
	return errores
}

func ValidarHorario(horario *obj.Valida_horario) ([]error) {
	var errores []error

	erroresMaterias := validarMaterias(horario.Materias, horario.Distribuciones)
	if len(erroresMaterias) > 0 {
		errores = append(errores, erroresMaterias...)
	}

	erroresBloques := validarBloques(horario.Salones, horario.Distribuciones)
	if len(erroresBloques) > 0 {
		errores = append(errores, erroresBloques...)
	}

	erroresProfesores := validarProfesores(horario.Profesores, horario.Distribuciones)
	if len(erroresProfesores) > 0 {
		errores = append(errores, erroresProfesores...)
	}

	return errores
}
