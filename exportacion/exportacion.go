package exportacion

import (
	"fmt"
	"strings"
	"os"
	"strconv"
	"sort"
	"time"
	"encoding/base64"
	"sync"
	utils "proyecto-horarios/utils_exportacion"
	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	obj "proyecto-horarios/objetos_exportacion"
)

var (
	contenidoTabla [][]string
	primerTitulo string
	segundoTitulo string
	tercerTitulo string
	idToNombre map[string]string
	clases map[string][]info_profe
	archivos []string
	errorGlobal error
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

type info_profe struct {
	Materia string
	Modulos []obj.Modulo
}

// Cosas del archivo HTML

func getEstilo() (string) {
	html := `
			<style>
			table, th, td {
				border: 1px solid black;
				border-collapse: collapse;
				text-align: center;
			}
			</style>
	`
	return html
}

func getLogos() (string) {
	html := `
			<img
					src="https://www.escom.ipn.mx/images/conocenos/escudoESCOM.png"
					style="width:150px;height:120px;"
			>
			<img
					src="https://upload.wikimedia.org/wikipedia/commons/f/f8/Logo_Instituto_Polit%C3%A9cnico_Nacional.png"
					align="right"
					style="width:250px;height:120px;"
			>
	`
	return html
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
	html := fmt.Sprintf("<td"+span+"><center>%s</center></td>", contenido)
	return html
}

func crearCeldaBloque(bloque obj.Bloque, span string) (string){
	html := "<td"+span+"><center>"
	html += fmt.Sprintf("<b> %s </b> <br>", bloque.Nombre)
	for i, m := range bloque.Modulos {
		html += fmt.Sprintf("%s %s - %s", m.Dia, m.Entrada, m.Salida)
		if i != len(bloque.Modulos) - 1 {
			html += "<br>"
		}
	}
	html += "</center></td>"
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

	html := `
	<html>
		<head>
		`+getEstilo()+`
		</head>
		<body>
		`+getLogos()+`
		<center>
	`
	// Agregamos titulo
	html += "<h1>" + "Horario" + "</h1>"
	// Iniciamos tabla
	html += "<table cellpadding = '5' cellspacing = '5'>"
	// Titulos de cada columna
	html += "<tr>"+
						"<th>"+primerTitulo+"</th>"+
						"<th>"+segundoTitulo+"</th>"+
						"<th>"+tercerTitulo+"</th>"+
					"</tr>"

	// Anadimos cada fila de la tabla
	for i, _ := range contenidoTabla {
		html += "<tr>"
		for _, cc := range contenidoTabla[i] {
			html += cc
		}
		html += "</tr>"
	}

	// Cerramos tabla
	html += "</table>"
	html += "</center></body></html>"

	return html
}

func exportarHorarioLista(horario *obj.Entrada_exportacion, ruta string) (string, error){
	tuplas := getTuplasIndividuales(horario)
	contenidoTabla = make([][]string, len(tuplas))

	if horario.Agrupar == "Profesor" {
		agruparProfesores(tuplas)
	} else if horario.Agrupar == "Materia" {
		agruparMaterias(tuplas)
	} else if horario.Agrupar == "Bloque" {
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


func crearContenidoHorarioIndividual(idProfe string) (string) {
	html := ""
	dias := [7]string{"Lunes", "Martes", "Miercoles", "Jueves", "Viernes", "Sabado", "Domingo"}
	for _, info := range clases[idProfe] {
		fila := `
			<tr>
				<td>`+info.Materia+`</td>
		`
		infoDias := make(map[string]string)
		for _, m := range info.Modulos {
			infoDias[m.Dia] = fmt.Sprintf("%s - %s", m.Entrada, m.Salida)
		}

		for _,d  := range dias {
			fila += "<td>"
			if cont, ok := infoDias[d]; ok {
				fila += cont
			} else {
				fila += " --- "
			}
			fila += "</td>"
		}
		fila += "</tr>"
		html += fila
	}
	return html
}

func crearTablaHorarioIndividual(idProfe string) (string) {
	nombreProfesor := idToNombre[idProfe]
	contenido := crearContenidoHorarioIndividual(idProfe)
	html := `
	<html>
		<head>
		 `+getEstilo()+`
		</head>
		<body>
			`+getLogos()+`
			<center>
				<h2>Horario</h2>
				<p> <b> Profesor/a: </b> `+nombreProfesor+` </p>
				<table cellpadding = '5' cellspacing = '5' style="width:75%">
					<tr>
						<th>Materia</th>
						<th>Lunes</th>
						<th>Martes</th>
						<th>Miercoles</th>
						<th>Jueves</th>
						<th>Viernes</th>
						<th>Sabado</th>
						<th>Domingo</th>
					</tr> `+contenido+`
				</table>
			</center>
		</body>
	</html>
	`

	return html
}

// Devuelve el nombre del archivo a comprimir
func crearHorarioIndividual(idProfe, ruta string) {
	// Modo landscape
	// Poner logos del IPN
	nombreProfe := idToNombre[idProfe]
	fecha := time.Now().Format("2006-01-02")
	nombreProfeFiltrado := strings.Replace(nombreProfe, " ", "-", -1)
	nombreArchivo := fmt.Sprintf("%s-%s.pdf",nombreProfeFiltrado, fecha)

	// Create new PDF generator
  pdfg, err := pdf.NewPDFGenerator()
  if err != nil {
		errorGlobal = err
		return
  }

	// Set options
	pdfg.Orientation.Set(pdf.OrientationLandscape)

	html := crearTablaHorarioIndividual(idProfe)

	// Add to document
	pdfg.AddPage(pdf.NewPageReader(strings.NewReader(html)))

  // Create PDF document in internal buffer
  err = pdfg.Create()
  if err != nil {
		errorGlobal = err
		return
  }

  // Write buffer contents to file on disk
  err = pdfg.WriteFile(ruta + nombreArchivo)
  if err != nil {
		errorGlobal = err
		return
  }


	archivos = append(archivos, ruta + nombreArchivo)
}


func exportarHorarioIndividual(horario *obj.Entrada_exportacion, ruta string) (string, error){
	idToNombre = make(map[string]string)
	clases = make(map[string][]info_profe)

	var profes []string

	for _, d := range horario.Distribuciones {
		modulos := d.Bloque.Modulos
		for _, a := range d.Asignaciones {
			if _, ok := idToNombre[a.Id_profesor]; !ok {
				idToNombre[a.Id_profesor] = a.Profesor
				profes = append(profes, a.Id_profesor)
			}

			c := info_profe{
				Materia: a.Materia,
				Modulos: modulos,
			}
			clases[a.Id_profesor] = append(clases[a.Id_profesor], c)
		}
	}

	sort.Slice(profes, func(i, j int) bool {
		return idToNombre[profes[i]] < idToNombre[profes[j]]
	})

	// Aqui va la parte asincrona con hilos
	var wg sync.WaitGroup

	for _, id := range profes {
		wg.Add(1)
		go func(id, ruta string) {
			defer wg.Done()
			crearHorarioIndividual(id, ruta)
		}(id, ruta)
	}
	wg.Wait()

	if errorGlobal != nil {
		return "", errorGlobal
	}

	// Aqui debemos hacer el zip de la carpeta tmp
	err := utils.ZipFiles(ruta+"horarios.zip", archivos)
	if err != nil {
		return "", err
	}

  content, err := os.ReadFile(ruta+"horarios.zip")
	if err != nil {
		return "", err
	}

	cadenaCodificada := base64.StdEncoding.EncodeToString(content)

	return cadenaCodificada, nil
}

func ExportarHorario(horario *obj.Entrada_exportacion, ruta string) (string, error){
	if horario.Tipo == "Lista" {
		return exportarHorarioLista(horario, ruta)
	} else if horario.Tipo == "Individual" {
		return exportarHorarioIndividual(horario, ruta)
	}

	return "", fmt.Errorf("%s no es un tipo valido de exportacion", horario.Tipo)
}
