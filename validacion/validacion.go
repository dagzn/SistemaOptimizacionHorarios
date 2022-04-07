package validacion

import (
	"fmt"
	obj "proyecto-horarios/objetos"
)

const (
	errorMaterias = "La materia '%s' tiene que impartirse %d veces y tiene %d clases asignadas."
	errorBloques = "El bloque %s tiene un maximo de %d salones y se le asignaron %d clases."
	errorClasesFaltantes = "El profesor %s requiere de %d clases en total y se le asignaron %d."
	errorInterseccionBloques = "El profesor %s tiene mas de una clase asignada para el bloque %s."
)

/*
	1.- Validar que cada materia sea dada una cierta cantidad de veces.
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
			errores = append(errores, fmt.Errorf(errorMaterias, m.Nombre, m.Cantidad, asignadas))
		}
	}

	return errores
}

/*
	1.- Validar que cada bloque cumpla con su limite asignado de salones.
*/
func validarBloques(salonesDisponibles int, distribuciones []obj.Distribucion) ([]error) {
	salonesOcupados := make(map[string]int)
	for _, d := range distribuciones {
		nombreBloque := d.Bloque
		actual := salonesOcupados[nombreBloque]
		salonesOcupados[nombreBloque] = actual + 1
	}

	var errores []error
	for _, b := range bloques {
		ocupados := salonesOcupados[b.Nombre]
		if salonesDisponibles < ocupados {
			errores = append(errores, fmt.Errorf(errorBloques, b.Nombre, salonesDisponibles, ocupados)
		}
	}

	return errores
}

/*
	1.- Validar que un profesor no tenga mas de una clase en el mismo bloque.
	2.- Validar que se cumplan con las clases requeridas para cada profesor.
	3.- Que el profesor de materias que si quiere dar. (PENDIENTE)
	4.- Que el profesor tenga asignados horarios en los que si pueda trabajar. (PENDIENTE)
*/
func validarProfesores(profesores []obj.Profesor, distribuciones []obj.Distribuciones) ([]error) {
	var errores []error

	clasesAsignadas := make(map[string]int)
	for _, d := range distribuciones {
		duplicados := make(map[string]bool)
		for _, a := range d.Asignaciones {
			profesor := a.Profesor
			if duplicados[profesor] {
				errores = append(errores, fmt.Errorf(errorInterseccionBloques, profesor, d.Bloque)
			}
			else {
				duplicados[profesor] = true
			}

			actual := clasesAsignadas[profesor]
			clasesAsignadas[profesor] = actual + 1
		}
	}

	for _, p := range profesores {
		asignadas := clasesAsignadas[p.Nombre]
		if asignadas != p.Clases {
			errores = append(errores, fmt.Errorf(errorClasesFaltantes, p.Nombre, p.Clases, asignadas))
		}
	}

	return errores
}

func ValidarHorario(horario *Valida_horario) ([]error) {
	var errores []error

	erroresMaterias := validarMaterias(horario.Materias, horario.Distribuciones)
	if len(erroresMaterias) > 0 {
		errores = append(errores, erroresMaterias)
	}

	erroresBloques := validarBloques(horario.Salones, horario.Distribuciones)
	if len(erroresBloques) > 0 {
		errores = append(errores, erroresBloques)
	}

	erroresProfesores := validarProfesores(horario.Profesores, horario.Distribuciones)
	if len(erroresProfesores) > 0 {
		errores = append(errores, erroresProfesores)
	}

	return errores
}
