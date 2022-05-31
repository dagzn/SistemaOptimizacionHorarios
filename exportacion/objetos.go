package exportacion

// Mantener los nombres de variables en singular

type Materia struct {
	Id       string `json:"id" validate:"required"`
	Nombre   string `json:"nombre" validate:"required"`
	Cantidad int    `json:"cantidad" validate:"required,gte=1"`
}

type Modulo struct {
	Dia     string `json:"dia"`
	Entrada string `json:"entrada"`
	Salida  string `json:"salida"`
}

type Bloque struct {
	Id      string   `json:"id" validate:"required"`
	Nombre  string   `json:"nombre" validate:"required"`
	Modulos []Modulo `json:"modulos"`
}

type Estructura struct {
	Id      string
	Bloques []Bloque
}

type Pref_materia struct {
	Id          string `json:"id" validate:"required"`
	Limite      int    `json:"limite" validate:"required"`
	Preferencia int    `json:"preferencia" validate:"required,oneof=-1 1 2 4 8 16 32 64 128 256 512"`
}

type Pref_bloque struct {
	Id          string `json:"id" validate:"required"`
	Preferencia int    `json:"preferencia" validate:"required,oneof=-1 1 2 4 8 16 32 64 128 256 512"`
}

type Profesor struct {
	Id       string         `json:"id" validate:"required"`
	Nombre   string         `json:"nombre" validate:"required"`
	Clases   int            `json:"clases" validate:"required,gte=1"`
	Materias []Pref_materia `json:"materias" validate:"required,min=1,dive,required"`
	Bloques  []Pref_bloque  `json:"bloques" validate:"required,min=1,dive,required"`
}

type Asignacion struct {
	Profesor    string `json:"profesor" validate:"required"`
	Id_profesor string `json:"id_profesor" validate:"required"`
	Materia     string `json:"materia" validate:"required"`
	Id_materia  string `json:"id_materia" validate:"required"`
}

type Distribucion struct {
	Bloque       Bloque       `json:"bloque" validate:"required"`
	Asignaciones []Asignacion `json:"asignaciones" validate:"required,min=1,dive,required"`
}

// Horario que exportaremos a PDF
type Entrada_exportacion struct {
	Distribuciones []Distribucion `validate:"required,min=1,dive,required"`
	Tipo           string         `validate:"required"` // Lista o Individual
	Agrupar        string         // Materia, Profesor, Bloque o NULL (sin agrupar)
}

type Salida_exportacion_fallida struct {
	Error string   `json:"error"`
	Logs  []string `json:"logs"`
}
