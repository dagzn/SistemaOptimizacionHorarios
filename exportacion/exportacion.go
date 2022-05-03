package exportacion

import (
	"fmt"
	"strings"
	"os"
	"strconv"
	"sort"
	"encoding/base64"
	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	obj "proyecto-horarios/objetos"
)

var (
	contenidoTabla [][]string
	primerTitulo string
	segundoTitulo string
	tercerTitulo string
)

type tupla struct {
	Bloque obj.Bloque
	Profesor string
	Materia string
}

type grupo struct {
	 Principal string
	 tamSpan int
}

type pareja struct {
	Primer  string
	Segundo string
}

func getTuplasIndividuales(horario *obj.Entrada_exportacion) ([]tupla) {
	var tuplas []tupla

	for _, d := range horario.Distribuciones {
		for _, a := range d.Asignaciones {
			t := tupla{
				Bloque: d.Bloque,
				Profesor: a.Profesor,
				Materia: a.Materia,
			}

			tuplas = append(tuplas, t)
		}
	}

	return tuplas
}

func horarioSencillo(tuplas []tupla) {
	for i, t := range tuplas {
		contenidoTabla[i] = append(contenidoTabla[i], crearCeldaSimple(t.Profesor, ""), crearCeldaSimple(t.Materia, ""), crearCeldaBloque(t.Bloque, ""))
	}

	primerTitulo, segundoTitulo, tercerTitulo = "Profesores", "Materias", "Bloques"
}

func crearCeldaSimple(contenido string, span string) (string){
	html := fmt.Sprintf("<td"+span+">%s</td>", contenido)
	return html
}

func crearCeldaBloque(bloque obj.Bloque, span string) (string){
	html := "<td"+span+">"
	html += fmt.Sprintf("<b> %s </b> <br>", bloque.Nombre)
	for i, m := range bloque.Modulos {
		html += fmt.Sprintf("%s %s - %s", m.Dia, m.Entrada, m.Salida)
		if i != len(bloque.Modulos) - 1 {
			html += "<br>"
		}
	}
	html += "</td>"
	return html
}

func agruparProfesores(tuplas []tupla) {
	sort.Slice(tuplas, func(i, j int) bool {
		if tuplas[i].Profesor == tuplas[j].Profesor {
			if tuplas[i].Materia == tuplas[j].Materia {
				return tuplas[i].Bloque.Nombre < tuplas[j].Bloque.Nombre
			}
			return tuplas[i].Materia < tuplas[j].Materia
		}
		return tuplas[i].Profesor < tuplas[j].Profesor
	})

	last := 0
	curr := 0
	for last < len(tuplas) {
		for curr < len(tuplas) && tuplas[curr].Profesor == tuplas[last].Profesor {
			curr++
		}
		rowspan := strconv.Itoa(curr - last)
		contenidoTabla[last] = append(contenidoTabla[last], crearCeldaSimple(tuplas[last].Profesor," rowspan='"+rowspan+"'"))
		last = curr
	}

	// llenamos con los demas
	for i, t := range tuplas {
		contenidoTabla[i] = append(contenidoTabla[i], crearCeldaSimple(t.Materia, ""), crearCeldaBloque(t.Bloque, ""))
	}

	primerTitulo, segundoTitulo, tercerTitulo = "Profesores", "Materias", "Bloques"
}

func agruparMaterias(tuplas []tupla) {
	sort.Slice(tuplas, func(i, j int) bool {
		if tuplas[i].Materia == tuplas[j].Materia {
			if tuplas[i].Profesor == tuplas[j].Profesor {
				return tuplas[i].Bloque.Nombre < tuplas[j].Bloque.Nombre
			}
			return tuplas[i].Profesor < tuplas[j].Profesor
		}
		return tuplas[i].Materia < tuplas[j].Materia
	})

	last := 0
	curr := 0
	for last < len(tuplas) {
		for curr < len(tuplas) && tuplas[curr].Materia == tuplas[last].Materia {
			curr++
		}
		rowspan := strconv.Itoa(curr - last)
		contenidoTabla[last] = append(contenidoTabla[last], crearCeldaSimple(tuplas[last].Materia, " rowspan='"+rowspan+"'"))
		last = curr
	}

	// llenamos con los demas
	for i, t := range tuplas {
		contenidoTabla[i] = append(contenidoTabla[i], crearCeldaSimple(t.Profesor, ""), crearCeldaBloque(t.Bloque, ""))
	}

	primerTitulo, segundoTitulo, tercerTitulo = "Materias", "Profesores", "Bloques"
}

func agruparBloques(tuplas []tupla) {
	sort.Slice(tuplas, func(i, j int) bool {
		if tuplas[i].Bloque.Nombre == tuplas[j].Bloque.Nombre {
			if tuplas[i].Profesor == tuplas[j].Profesor {
				return tuplas[i].Materia < tuplas[j].Materia
			}
			return tuplas[i].Profesor < tuplas[j].Profesor
		}
		return tuplas[i].Bloque.Nombre < tuplas[j].Bloque.Nombre
	})

	last := 0
	curr := 0
	for last < len(tuplas) {
		for curr < len(tuplas) && tuplas[curr].Bloque.Id == tuplas[last].Bloque.Id {
			curr++
		}
		rowspan := strconv.Itoa(curr - last)
		contenidoTabla[last] = append(contenidoTabla[last], crearCeldaBloque(tuplas[last].Bloque, " rowspan='"+rowspan+"'"))
		last = curr
	}

	// llenamos con los demas
	for i, t := range tuplas {
		contenidoTabla[i] = append(contenidoTabla[i], crearCeldaSimple(t.Profesor, ""), crearCeldaSimple(t.Materia, ""))
	}

	primerTitulo, segundoTitulo, tercerTitulo = "Bloques", "Profesores", "Materias"
}

func crearTabla() (string) {

	html := "<html><body><center>"
	// Agregamos titulo
	html += "<h1>" + "Horario" + "</h1>"
	// Iniciamos tabla
	html += "<table border = '1' cellpadding = '5' cellspacing = '5'>"
	// Titulos de cada columna
	html += "<tr>"+
						"<th>"+primerTitulo+"</th>"+
						"<th>"+segundoTitulo+"</th>"+
						"<th>"+tercerTitulo+"</th>"+
					"</tr>"

	// Anadimos cada fila de la tabla
	for i, _ := range contenidoTabla {
		html += "<tr>"
		fmt.Println("Linea numero ", i)
		for _, cc := range contenidoTabla[i] {
			html += cc
			fmt.Println(cc)
		}
		html += "</tr>"
	}

	// Cerramos tabla
	html += "</table>"
	html += "</center></body></html>"

	return html
}

func ExportarHorario(horario *obj.Entrada_exportacion, ruta string) (string, error){
	tuplas := getTuplasIndividuales(horario)
	contenidoTabla = make([][]string, len(tuplas))

	if horario.Agrupar == "Profesores" {
		agruparProfesores(tuplas)
	} else if horario.Agrupar == "Materias" {
		agruparMaterias(tuplas)
	} else if horario.Agrupar == "Bloques" {
	 agruparBloques(tuplas)
	} else {
		horarioSencillo(tuplas)
	}

	html := crearTabla()

	// Create new PDF generator
  pdfg, err := pdf.NewPDFGenerator()
  if err != nil {
		return "", err
  }

	// Add to document
	pdfg.AddPage(pdf.NewPageReader(strings.NewReader(html)))

  // Create PDF document in internal buffer
  err = pdfg.Create()
  if err != nil {
		return "", err
  }

  // Write buffer contents to file on disk
  err = pdfg.WriteFile(ruta + "simplesample.pdf")
  if err != nil {
		return "", err
  }

  content, err := os.ReadFile(ruta + "simplesample.pdf")
	if err != nil {
		return "", err
	}

	cadenaCodificada := base64.StdEncoding.EncodeToString(content)

	return cadenaCodificada, nil
}

