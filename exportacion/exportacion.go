package exportacion

import (
	"fmt"
	"strings"
	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	obj "proyecto-horarios/objetos"
)

type tupla struct {
	Bloque obj.Bloque
	Profesor string
	Materia string
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

func crearHTMLBasico(horario *obj.Entrada_exportacion) (string) {
	tuplas := getTuplasIndividuales(horario)

	html := "<html><body><center>"
	// Agregamos titulo
	html += "<h1>" + "Horario" + "</h1>"
	// Iniciamos tabla
	html += "<table border = '1' cellpadding = '5' cellspacing = '5'>"
	// Titulos de cada columna
	html += "<tr>"+
						"<th>Profesor</th>"+
						"<th>Materia</th>"+
						"<th>Horario</th>"+
					"</tr>"

	// Anadimos cada fila de la tabla
	for _, t := range tuplas {
		html += "<tr>"
		// Nombre del profesor
		html += fmt.Sprintf("<td>%s</td>", t.Profesor)
		// Nombre de la materia
		html += fmt.Sprintf("<td>%s</td>", t.Materia)
		// Horarios
		html += "<td>"
		for i, m := range t.Bloque.Modulos {
			html += fmt.Sprintf("%s %s - %s", m.Dia, m.Entrada, m.Salida)
			if i != len(t.Bloque.Modulos) - 1 {
				html += "<br>"
			}
		}
		html += "</td>"
		html += "</tr>"
	}

	// Cerramos tabla
	html += "</table>"
	html += "</center></body></html>"

	return html
}

func ExportarHorario(horario *obj.Entrada_exportacion) (error){
	html := crearHTMLBasico(horario)

	// Create new PDF generator
  pdfg, err := pdf.NewPDFGenerator()
  if err != nil {
		return err
  }

	// Add to document
	pdfg.AddPage(pdf.NewPageReader(strings.NewReader(html)))

  // Create PDF document in internal buffer
  err = pdfg.Create()
  if err != nil {
		return err
  }

  // Write buffer contents to file on disk
  err = pdfg.WriteFile("./simplesample.pdf")
  if err != nil {
		return err
  }

  fmt.Println("Done")

	return nil
}

