package objetos

// Mantener los nombres de variables en singular

type Materia struct {
	Id       string `validate:"required"`
	Nombre   string `validate:"required"`
	Cantidad int    `validate:"required,gte=1"`
}

type Modulo struct {
	Dia     string
	Entrada string
	Salida  string
}

type Bloque struct {
	Id      string `validate:"required"`
	Nombre  string `validate:"required"`
	Modulos []Modulo
}

type Estructura struct {
	Id      string
	Bloques []Bloque
}

type Pref_materia struct {
	Id          string `validate:"required"`
	Limite      int    `validate:"required"`
	Preferencia int    `validate:"required,oneof=-1 1 2 4 8 16 32 64 128 256 512"`
}

type Pref_bloque struct {
	Id          string `validate:"required"`
	Preferencia int    `validate:"required,oneof=-1 1 2 4 8 16 32 64 128 256 512"`
}

type Profesor struct {
	Id       string         `validate:"required"`
	Nombre   string         `validate:"required"`
	Clases   int            `validate:"required,gte=1"`
	Materias []Pref_materia `validate:"required,min=1,dive,required"`
	Bloques  []Pref_bloque  `validate:"required,min=1,dive,required"`
}

type Asignacion struct {
	Profesor    string `validate:"required"`
	Id_profesor string `validate:"required"`
	Materia     string `validate:"required"`
	Id_materia  string `validate:"required"`
}

type Distribucion struct {
	Bloque       Bloque       `validate:"required"`
	Asignaciones []Asignacion `validate:"required,min=1,dive,required"`
}

// Formato de entrada para crear un horario
type Entrada_horario struct {
	Salones    int        `validate:"required,gte=1"`
	Materias   []Materia  `validate:"required,min=1,dive,required"`
	Profesores []Profesor `validate:"required,min=1,dive,required"`
	Bloques    []Bloque   `validate:"required,min=1,dive,required"`
}

// Formato de salida al crear un horario
type Salida_horario struct {
	Distribuciones []Distribucion
	Error          string
	Logs           []string
}

// Este formato es el que usamos para validar
type Entrada_validacion struct {
	Distribuciones []Distribucion `validate:"required,min=1,dive,required"`
	Profesores     []Profesor     `validate:"required,min=1,dive,required"`
	Materias       []Materia      `validate:"required,min=1,dive,required"`
	Salones        int            `validate:"required,gte=1"`
	Validaciones   []string       `validate:"required,min=1,dive,required"`
}

// Coleccion de los errores encontrados al validar
type Salida_validacion struct {
	Error string
	Logs  []string
}

// Horario que exportaremos a PDF
type Entrada_exportacion struct {
	Distribuciones []Distribucion `validate:"required,min=1,dive,required"`
	Agrupar        string         // Materia, Profesor, Bloque o NULL (sin agrupar)
}
