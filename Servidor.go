package main

import (
	"container/list"
	"errors"
	"fmt"
	"net"
	"net/rpc"
)

const conexion = ":9999"
const protocoloConexion = "tcp"

type Servidor struct {
	Admin Administrador
}

func (serv *Servidor) AgregarCalificacionAlumno(mensaje Mensaje, reply *string) error {
	nombreMateria := mensaje.Materia
	nombreAlumno := mensaje.Alumno
	calificacion := mensaje.Calificacion
	
	if serv.Admin.existeMateria(nombreMateria) {
		if serv.Admin.existeAlumnoMateria(nombreAlumno, nombreMateria) {
			return errors.New("El alumno ya tiene registrada una calificación en esta materia")
		}
		serv.Admin.agregarAlumno(nombreAlumno, nombreMateria, calificacion)
		*reply = "\n" + "Ha sido agregado a la materia, un nuevo alumno:" + " (" + nombreAlumno + ")"
	} else {
		serv.Admin.agregarMateria(nombreMateria)
		serv.Admin.agregarAlumno(nombreAlumno, nombreMateria, calificacion)
		*reply = "\n" + "Una nueva materia ha sido agregada"+ "\n" + "Un nuevo alumno ha sido agregado a la materia"
	}
	*reply += "\n(" + nombreAlumno + " en la materia: " + nombreMateria + ", calificación = " + fmt.Sprintf("%f", calificacion) + ")"
	return nil
}

func (serv *Servidor) ObtenerPromedioAlum(nombreAlumno string, reply *float64) error {
	nombreAlumno = nombreAlumno
	if serv.Admin.existeAlumno(nombreAlumno) {
		*reply = serv.Admin.ObtenerPromedioAlum(nombreAlumno)
	} else {
		return errors.New("No existe el alumno: " + nombreAlumno)
	}
	return nil
}

func (serv *Servidor) ObtenerPromedioGeneralAlum(consulta string, reply *float64) error {
	if serv.Admin.Materias.Len() > 0 {
		*reply = serv.Admin.ObtenerPromedioGeneralAlum()
	} else {
		return errors.New("No hay materias registradas")
	}
	return nil
}

func (serv *Servidor) ObtenerPromedioMateria(nombreMateria string, reply *float64) error {
	nombreMateria = nombreMateria
	if serv.Admin.existeMateria(nombreMateria) {
		*reply = serv.Admin.ObtenerPromedioMateria(nombreMateria)
	} else {
		return errors.New("No existe la materia: " + nombreMateria)
	}
	return nil
}

func (serv *Servidor) ObtenerCalificacionesAlumno(nombreAlumno string, reply *string) error {
	nombreAlumno = nombreAlumno
	if serv.Admin.existeAlumno(nombreAlumno) {
		*reply = "\n" + serv.Admin.ObtenerCalificacionesAlumno(nombreAlumno)
	} else {
		return errors.New("No existe el alumno: " + nombreAlumno)
	}
	return nil
}

func (serv *Servidor) ObtenerCalificacionesMaterias(consulta string, reply *string) error {
	if serv.Admin.Materias.Len() > 0 {
		*reply = "\n" + serv.Admin.ObtenerCalificacionesMaterias()
	} else {
		return errors.New("No hay ninguna materia registrada")
	}
	return nil
}

func (serv *Servidor) ObtenerCalificacionesAlumnosMateria(nombreMateria string, reply *string) error {
	nombreMateria = nombreMateria
	if serv.Admin.existeMateria(nombreMateria) {
		*reply = "\n" + serv.Admin.ObtenerCalificacionesAlumnosMateria(nombreMateria)
	} else {
		return errors.New("No existe la materia: " + nombreMateria)
	}
	return nil
}

type Administrador struct {
	Materias list.List
}

func (admin *Administrador) existeMateria(nombreMateria string) bool {
	for e := admin.Materias.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Materia)
		if m.Nombre == nombreMateria {
			return true
		}
	}
	return false
}

func (admin *Administrador) obtenerMateria(nombreMateria string) *Materia {
	for e := admin.Materias.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Materia)
		if m.Nombre == nombreMateria {
			return m
		}
	}
	return nil
}

func (admin *Administrador) existeAlumnoMateria(nombreAlumno string, nombreMateria string) bool {
	m := admin.obtenerMateria(nombreMateria)
	return m.existeAlumno(nombreAlumno)
}

func (admin *Administrador) existeAlumno(nombreAlumno string) bool {
	for e := admin.Materias.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Materia)
		if m.existeAlumno(nombreAlumno) {
			return true
		}
	}
	return false
}

func (admin *Administrador) agregarAlumno(nombreAlumno string, nombreMateria string, calificacion float64) {
	m := admin.obtenerMateria(nombreMateria)
	m.agregarAlumno(nombreAlumno, calificacion)
}

func (admin *Administrador) agregarMateria(nombreMateria string) {
	m := new(Materia)
	m.Nombre = nombreMateria
	admin.Materias.PushBack(m)
}

func (admin *Administrador) ObtenerPromedioAlum(nombreAlumno string) float64 {
	cantMaterias := 0.0
	sumaCalificaciones := 0.0
	for e := admin.Materias.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Materia)
		if m.existeAlumno(nombreAlumno) {
			cantMaterias++
			sumaCalificaciones += m.obtenerAlumno(nombreAlumno).Calificacion
		}
	}
	return sumaCalificaciones / cantMaterias
}

func (admin *Administrador) ObtenerCalificacionesAlumno(nombreAlumno string) string {
	calificaciones := ""
	for e := admin.Materias.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Materia)
		if m.existeAlumno(nombreAlumno) {
			calificaciones += "* Materia:" + m.Nombre + ", calificación= " + fmt.Sprintf("%f", m.obtenerAlumno(nombreAlumno).Calificacion) + "\n"
		}
	}
	return calificaciones
}

func (admin *Administrador) ObtenerCalificacionesMaterias() string {
	calificaciones := ""
	for e := admin.Materias.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Materia)
		calificaciones += "* Materia: " + m.Nombre + ", promedio = " + fmt.Sprintf("%f", m.obtenerPromedio()) + "\n"
	}
	return calificaciones
}

func (admin *Administrador) ObtenerCalificacionesAlumnosMateria(nombreMateria string) string {
	m := admin.obtenerMateria(nombreMateria)
	return m.obtenerCalificacionesAlumnos()
}

func (admin *Administrador) ObtenerPromedioMateria(nombreMateria string) float64 {
	m := admin.obtenerMateria(nombreMateria)
	return m.obtenerPromedio()
}

func (admin *Administrador) ObtenerPromedioGeneralAlum() float64 {
	cantMaterias := 0.0
	sumaPromediosMaterias := 0.0
	for e := admin.Materias.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Materia)
		cantMaterias++
		sumaPromediosMaterias += m.obtenerPromedio()
	}
	return sumaPromediosMaterias / cantMaterias
}

type Materia struct {
	Nombre  string
	Alumnos list.List
}

func (m *Materia) existeAlumno(nombreAlumno string) bool {
	for e := m.Alumnos.Front(); e != nil; e = e.Next() {
		a := e.Value.(*Alumno)
		if a.Nombre == nombreAlumno {
			return true
		}
	}
	return false
}

func (m *Materia) obtenerAlumno(nombreAlumno string) *Alumno {
	for e := m.Alumnos.Front(); e != nil; e = e.Next() {
		a := e.Value.(*Alumno)
		if a.Nombre == nombreAlumno {
			return a
		}
	}
	return nil
}

func (m *Materia) modificarCalificacionAlumno(nombreAlumno string, calificacion float64) {
	a := m.obtenerAlumno(nombreAlumno)
	a.Calificacion = calificacion
}

func (m *Materia) agregarAlumno(nombreAlumno string, calificacion float64) {
	a := new(Alumno)
	a.Calificacion = calificacion
	a.Nombre = nombreAlumno
	m.Alumnos.PushBack(a)
}

func (m *Materia) obtenerPromedio() float64 {
	contAlumnos := 0.0
	sumaCalificaciones := 0.0
	for e := m.Alumnos.Front(); e != nil; e = e.Next() {
		a := e.Value.(*Alumno)
		contAlumnos++
		sumaCalificaciones += a.Calificacion
	}
	return sumaCalificaciones / contAlumnos
}

func (m *Materia) obtenerCalificacionesAlumnos() string {
	calificaciones := ""
	for e := m.Alumnos.Front(); e != nil; e = e.Next() {
		a := e.Value.(*Alumno)
		calificaciones += "* Alumno: " + a.Nombre + ", calificación = " + fmt.Sprintf("%f", m.obtenerAlumno(a.Nombre).Calificacion) + "\n"
	}
	return calificaciones
}

type Alumno struct {
	Nombre       string
	Calificacion float64
}

type Mensaje struct {
	Materia      string
	Alumno       string
	Calificacion float64
}

func servidor() {
	rpc.Register(new(Servidor))
	ln, err := net.Listen(protocoloConexion, conexion)
	if err != nil {
		fmt.Println(err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(c)
	}
}

func main() {
	var input string
	go servidor()
	fmt.Println("\n\n\tEl servidor se ha iniciado")
	fmt.Printf("\n\n\nPresione enter para terminar...")
	fmt.Scanln(&input)
	fmt.Printf("\n\nTerminar procesos del servidor")
}