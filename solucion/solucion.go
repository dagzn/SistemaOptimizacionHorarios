package solucion

import (
	"container/heap"
	"fmt"
)

const (
	inf                       = int64(1e18)
	oo                        = int64(1e9)
	errorClasesMaterias       = "Se determinaron %d clase/s a impartir, pero en total los profesores seleccionados deben impartir %d clase/s."
	errorClasesProfesores     = "En total los profesores deben dar %d clase/s, pero solo existe/n %d clase/s disponible/s."
	errorSinSolucion          = "No existe alguna asignación de horario válida dados los conjuntos de profesores, materias y bloques seleccionados."
	errorEspacioInsuficiente  = "Se desea impartir %d clase/s pero solamente hay %d bloque/s de horario disponibles con %d salon/es cada uno."
	errorBloquesInsuficientes = "El/la profesor/a %s solamente tiene %d horario/s disponible/s pero se le pretende asignar %d clase/s."
	errorProfesInsuficientes  = "La materia %s debe darse en %d clase/s pero, en total, los profesores seleccionados solamente pueden impartirla %d veces"
	logCostoInfinito          = "Solo fue posible encontrar una solución que hace uso de preferencias denominadas como no óptimas."
	logAristaInfinita         = "En el horario se usa la asignación de el/la profesor/a %s con %s %s, la cual se registró con una preferencia no óptima"
)

type tupla struct {
	Profesor Profesor
	Materia  Materia
	Bloque   Bloque
}

type edge struct {
	src, dst  int
	cap, flow int64
	cost      int64
	rev       int
}

type bedge struct {
	node int
	pos  int
}

var (
	n            int
	adj          [][]edge
	dist         []int64
	potential    []int64
	back         []bedge
	idx_original map[int]int
	logs         []string
)

// Dijkstra code
type Item struct {
	node     int
	distance int64
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(v interface{}) {
	*pq = append(*pq, v.(Item))
}

func (pq *PriorityQueue) Pop() (v interface{}) {
	a := *pq
	*pq, v = a[:len(a)-1], a[len(a)-1]
	return
}

// END Dijkstra code

func add_edge(src, dst int, cap, cost int64) {
	if cost == -1 {
		cost = oo
	}
	adj[src] = append(adj[src], edge{src, dst, cap, 0, cost, len(adj[dst])})
	adj[dst] = append(adj[dst], edge{dst, src, 0, 0, -cost, len(adj[src]) - 1})
}

func rcost(e edge) int64 {
	return e.cost + potential[e.src] - potential[e.dst]
}

func dijkstra(source, sink int) bool {

	dist = make([]int64, n+1)
	back = make([]bedge, n+1)
	for i := range dist {
		dist[i] = inf
		back[i] = bedge{-1, -1}
	}

	dist[source] = 0
	pq := PriorityQueue{Item{source, 0}}
	heap.Init(&pq)

	for len(pq) > 0 {
		curr := heap.Pop(&pq).(Item)
		u, d := curr.node, curr.distance
		if d != dist[u] {
			continue
		}
		for i, e := range adj[u] {
			if new_d := dist[e.src] + rcost(e); e.flow < e.cap && dist[e.dst] > new_d {
				dist[e.dst] = new_d
				back[e.dst] = bedge{u, i}
				heap.Push(&pq, Item{e.dst, dist[e.dst]})
			}
		}
	}

	return dist[sink] < inf
}

func min_cost_max_flow(source, sink int) (int64, int64) {
	var cost, flow int64
	potential = make([]int64, n+1)
	for dijkstra(source, sink) {
		for u := 0; u < n; u++ {
			if dist[u] < inf {
				potential[u] = dist[u]
			}
		}

		new_flow := inf
		for be := back[sink]; be.node != -1; be = back[be.node] {
			e := adj[be.node][be.pos]
			if e.cap-e.flow < new_flow {
				new_flow = e.cap - e.flow
			}
		}

		flow += new_flow
		for be := back[sink]; be.node != -1; be = back[be.node] {
			adj[be.node][be.pos].flow += new_flow
			aux := adj[be.node][be.pos]
			adj[aux.dst][aux.rev].flow -= new_flow
			cost += new_flow * aux.cost
		}

	}
	return flow, cost
}

func crearGrafo(salones int, materias []Materia, profesores []Profesor, bloques []Bloque) (int, int, error) {
	n_materias := len(materias)
	n_profesores := len(profesores)
	n_bloques := len(bloques)
	n = 1 + n_materias + 2*n_profesores + n_bloques + 1
	fuente := 0
	destino := n - 1
	clasesMaterias := 0
	clasesProfesores := 0

	adj = make([][]edge, n)
	m_idx := make(map[string]int)
	b_idx := make(map[string]int)
	sumaCapMaterias := make(map[string]int)
	idx_original = make(map[int]int)

	// indices para materias
	// aristas de fuente hacia materias
	for i, m := range materias {
		m_idx[m.Id] = i + 1
		idx_original[i+1] = i
		clasesMaterias += m.Cantidad
		add_edge(fuente, i+1, int64(m.Cantidad), 0)
	}

	// indices para bloques
	// aristas de bloques hacia destino
	for i, b := range bloques {
		nodo := 1 + n_materias + 2*n_profesores + i
		b_idx[b.Id] = nodo
		idx_original[nodo] = i
		add_edge(nodo, destino, int64(salones), 0)
	}

	for i, p := range profesores {
		entrada := 1 + n_materias + 2*i
		salida := 1 + n_materias + 2*i + 1
		idx_original[entrada] = i
		add_edge(entrada, salida, int64(p.Clases), 0)
		clasesProfesores += p.Clases

		pref_materias := p.Materias
		for _, m := range pref_materias {
			if nodo, ok := m_idx[m.Id]; ok {
				sumaCapMaterias[m.Id] += m.Limite
				add_edge(nodo, entrada, int64(m.Limite), int64(m.Preferencia))
			}
		}

		pref_bloques := p.Bloques
		for _, b := range pref_bloques {
			if nodo, ok := b_idx[b.Id]; ok {
				add_edge(salida, nodo, 1, int64(b.Preferencia))
			}
		}

		if cuantosBloques := len(pref_bloques); cuantosBloques < p.Clases {
			return 0, 0, fmt.Errorf(errorBloquesInsuficientes, p.Nombre, cuantosBloques, p.Clases)
		}
	}

	for _, m := range materias {
		if maxCap, ok := sumaCapMaterias[m.Id]; ok {
			if maxCap < m.Cantidad {
				return 0, 0, fmt.Errorf(errorProfesInsuficientes, m.Nombre, m.Cantidad, maxCap)
			}
		}
	}

	if clasesMaterias > clasesProfesores {
		// Se planea dar mas clases de las que los profesores deben dar
		return 0, 0, fmt.Errorf(errorClasesMaterias, clasesMaterias, clasesProfesores)
	} else if clasesMaterias < clasesProfesores {
		// Los profesores necesitan mas clases de las disponibles
		return 0, 0, fmt.Errorf(errorClasesProfesores, clasesProfesores, clasesMaterias)
	} else if espacio := salones * len(bloques); espacio < clasesProfesores{
		// En total no hay suficientes salones para todas las clases.
		return 0, 0, fmt.Errorf(errorEspacioInsuficiente, clasesProfesores, len(bloques), salones)
	}

	return fuente, destino, nil
}

func filtrarTuplasPorBloque(tuplas []tupla, id_bloque string) []Asignacion {
	var asignaciones []Asignacion

	for _, t := range tuplas {
		if t.Bloque.Id == id_bloque {
			a := Asignacion{
				Profesor:    t.Profesor.Nombre,
				Id_profesor: t.Profesor.Id,
				Materia:     t.Materia.Nombre,
				Id_materia:  t.Materia.Id,
			}
			asignaciones = append(asignaciones, a)
		}
	}

	return asignaciones
}

func movimiento(u int) (int, bool) {
	for i, e := range adj[u] {
		if e.cost < 0 {
			continue
		}
		if e.flow <= 0 {
			continue
		}
		adj[u][i].flow--
		return e.dst, (adj[u][i].cost >= oo)
	}
	return -1, false
}

func encontrarSolucion(fuente, destino int, materias []Materia, profesores []Profesor, bloques []Bloque) ([]Distribucion, error) {
	flujo, costo := min_cost_max_flow(fuente, destino)

	// Aqui es donde reconstruimos la respuesta
	var distribuciones []Distribucion
	var tuplas []tupla
	var flujoEsperado int64

	for _, e := range adj[fuente] {
		flujoEsperado += e.cap
	}

	if flujo != flujoEsperado {
		return nil, fmt.Errorf(errorSinSolucion)
	}

	if costo >= oo {
		// Aqui debemos indicar cuales preferencias se pueden modificar.
		// Esto consiste en identificar las aristas de costo infinito que tengan flujo.
		logs = append(logs, fmt.Sprintf(logCostoInfinito))
	}

	flujoFinal := flujo
	for i := int64(0); i < flujoFinal; i++ {
		u := 0
		u, _ = movimiento(u)
		materia := materias[idx_original[u]]

		u, infinita := movimiento(u)
		profesor := profesores[idx_original[u]]
		if infinita {
			logs = append(logs, fmt.Sprintf(logAristaInfinita, profesor.Nombre,"la materia", materia.Nombre))
		}

		u, _ = movimiento(u)
		u, infinita = movimiento(u)
		bloque := bloques[idx_original[u]]
		if infinita {
			logs = append(logs, fmt.Sprintf(logAristaInfinita, profesor.Nombre,"el bloque", bloque.Nombre))
		}

		// Optimizacion: agregar directo a un map de bloque -> {profesor,materia}
		t := tupla{
			Profesor: profesor,
			Materia:  materia,
			Bloque:   bloque,
		}

		tuplas = append(tuplas, t)
	}

	for _, b := range bloques {
		asignaciones := filtrarTuplasPorBloque(tuplas, b.Id)
		if len(asignaciones) > 0 {
			d := Distribucion{
				Bloque:       b,
				Asignaciones: asignaciones,
			}

			distribuciones = append(distribuciones, d)
		}
	}

	return distribuciones, nil
}

func GenerarHorario(horario *Entrada_horario) (*Salida_horario, error) {
	fuente, destino, err := crearGrafo(horario.Salones, horario.Materias, horario.Profesores, horario.Bloques)
	if err != nil {
		return nil, err
	}

	distribuciones, err := encontrarSolucion(fuente, destino, horario.Materias, horario.Profesores, horario.Bloques)
	if err != nil {
		// TODO: Aqui regresar salida con error y logs si encontramos como mejorar el no hay resultado.
		return nil, err
	}

	// Aqui es donde creamos el obj Salida_horario y lo regresamos
	salida := &Salida_horario{
		Distribuciones: distribuciones,
		Logs:           logs,
	}

	return salida, nil
}
