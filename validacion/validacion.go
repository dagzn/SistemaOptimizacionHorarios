package validacion

import (
	"fmt"
)

const (
	errorValidacion = "Se encontraron errores al momento de validar el horario."
	errorClasesMateria = "La materia '%s' tiene que impartirse %d veces y en el horario tiene %d clases asignadas."
	errorClasesProfesor = "El/la profesor/a %s requiere de %d clase/s en total y en el horario tiene %d clase/s."
	errorInterseccionBloques = "El/la profesor/a %s tiene más de una clase asignada para el bloque %s."
	errorLimiteBloque = "El bloque %s tiene un máximo de %d salones y se le asignaron %d clases."
	errorLimiteMateria = "El limite impuesto para el/la profesor/a %s para dar la materia %s es de %d clases, pero se le asignaron %d clase/s."
	errorMateriaNoAsignada = "La materia %s no puede ser impartida por el/la profesor/a %s."
	errorBloqueNoAsignado = "El/la profesor/a %s no puede impartir clases en el bloque %s."
)

var (
	materiasAsignadas map[string][]string
	bloquesAsignados map[string][]string
	nombreMateria map[string]string
	nombreBloque map[string]string
)

/*
	Validar que cada materia sea dada una cierta cantidad de veces.
*/
func validarMaterias(materias []Materia, distribuciones []Distribucion, validaciones map[string]int) ([]error) {
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
func validarBloques(salonesDisponibles int, distribuciones []Distribucion, validaciones map[string]int) ([]error) {
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

func unique(stringSlice []string) []string {
    keys := make(map[string]bool)
    list := []string{}
    for _, entry := range stringSlice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
}

// Validar que un profesor no tenga mas de una clase en el mismo bloque.
func validarInterseccionBloques(distribuciones []Distribucion, validaciones map[string]int) ([]error) {
	if _, ok := validaciones["InterseccionBloques"]; !ok {
		return nil
	}

	var errores []error
	for _, d := range distribuciones {
		duplicados := make(map[string]bool)
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
func validarClasesProfesor(profesores []Profesor, validaciones map[string]int) ([]error) {
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
func validarLimiteMateria(profesores []Profesor, validaciones map[string]int) ([]error) {
	if _, ok := validaciones["LimiteMateria"]; !ok {
		return nil
	}

	var errores []error
	for _, p := range profesores {
		vecesAsignada := make(map[string]int)
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
func validarMateriaNoAsignada(profesores []Profesor, validaciones map[string]int) ([]error) {
	if _, ok := validaciones["MateriaNoAsignada"]; !ok {
		return nil
	}

	var errores []error
	for _, p := range profesores {
		preferencias := make(map[string]int)
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
func validarBloqueNoAsignado(profesores []Profesor, validaciones map[string]int) ([]error) {
	if _, ok := validaciones["BloqueNoAsignado"]; !ok {
		return nil
	}

	var errores []error
	for _, p := range profesores {
		preferencias := make(map[string]int)
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

func validarProfesores(profesores []Profesor, distribuciones []Distribucion, validaciones map[string]int) ([]error) {
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

func llenarInformacionProfesores(distribuciones []Distribucion) {
	materiasAsignadas = make(map[string][]string)
	bloquesAsignados = make(map[string][]string)
	nombreMateria = make(map[string]string)
	nombreBloque = make(map[string]string)

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

func validar(horario *Entrada_validacion) ([]error) {
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

func ValidarHorario(horario *Entrada_validacion) (*Salida_validacion){
	errores := validar(horario)
	if len(errores) > 0 {
		err := fmt.Errorf(errorValidacion);
		salida := &Salida_validacion{
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

	return &Salida_validacion{}
}
