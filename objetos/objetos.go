package objetos

// Mantener los nombres de variables en singular

type Materia struct {
	Id int
	Nombre string
	Cantidad int
}

type Modulo struct {
	Dia string
	Entrada string
	Salida string
}

type Bloque struct {
	Id int
	Nombre string
	Modulos []Modulo
}

type Estructura struct {
	Id int
	Bloques []Bloque
}

type Pref_materia struct {
	Id int
	Limite int
	Preferencia int
}

type Pref_bloque struct {
	Id int
	Preferencia int
}

type Profesor struct {
	Id int
	Nombre string
	Clases int
	Materias []Pref_materia
	Bloques []Pref_bloque
}

type Asignacion struct {
	Profesor string
	Id_profesor int
	Materia string
	Id_materia int
}

type Distribucion struct {
	Bloque string
	Id_bloque int
	Asignaciones []Asignacion
}

// Formato de entrada para crear un horario
type Entrada_horario struct {
	Salones int
	Materias []Materia
	Profesores []Profesor
	Bloques []Bloque
}

// Formato de salida al crear un horario
type Salida_horario struct {
	Distribuciones []Distribucion
	Profesores []Profesor
	Materias []Materia
}

// Este formato es el que usamos para validar
type Valida_horario struct {
	Distribuciones []Distribucion
	Profesores []Profesor
	Materias []Materia
	Salones int
}
