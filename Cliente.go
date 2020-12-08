package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
)

const conexion = "127.0.0.1:9999"
const protocoloConexion = "tcp"

//Nombres para los metodos del Servidor
const nombreServidor = "Servidor"
const metodoAgregarCalificacionAlumno = ".AgregarCalificacionAlumno"
const metodoPromedioAlumno = ".ObtenerPromedioAlum"
const metodoPromedioGeneral = ".ObtenerPromedioGeneralAlum"
const metodoPromedioMateria = ".ObtenerPromedioMateria"
const metodoCalificacionesAlumno = ".ObtenerCalificacionesAlumno"
const metodoCalifiacionesMaterias = ".ObtenerCalificacionesMaterias"
const metodoCalificacionesAlumnosMateria = ".ObtenerCalificacionesAlumnosMateria"

type Mensaje struct {
	Materia      string
	Alumno       string
	Calificacion float64
}

func main() {
	ejecutar()
}

func ejecutar() {
	var respuesta string
	var promedio float64
	continuar := true

	c, err := rpc.Dial(protocoloConexion, conexion)
	if err != nil {
		fmt.Println(err)
		return
	}
	for continuar {
		var opc int
		fmt.Println("\n\n**********Menu**********")
		fmt.Println("1) Agregar la calificación de una materia para un alumno")
		fmt.Println("2) Obtener promedio de un alumno")
		fmt.Println("3) Obtener promedio general de los alumnos")
		fmt.Println("4) Obtener promedio de una materia")
		fmt.Println("0) Salir")
		fmt.Print("Ingrese una opción: ")
		fmt.Scanln(&opc)

		switch opc {
		case 1:
			err = c.Call(nombreServidor+metodoAgregarCalificacionAlumno, agregarCalificacionAlumno(), &respuesta)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(respuesta)
			}
		case 2:
			fmt.Println("\n\t ---Obtener promedio de un alumno---")
			nombreAlumno := dameNombreAlumno()
			err = c.Call(nombreServidor+metodoPromedioAlumno, nombreAlumno, &promedio)
			if err != nil {
				fmt.Println(err)
			} else {
				err = c.Call(nombreServidor+metodoCalificacionesAlumno, nombreAlumno, &respuesta)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(respuesta)
					fmt.Println("Promedio del alumno = " + fmt.Sprintf("%f", promedio))
				}
			}
		case 3:
			fmt.Println("\t---Promedio general de todos los alumnos---")
			err = c.Call(nombreServidor+metodoPromedioGeneral, "", &promedio)
			if err != nil {
				fmt.Println(err)
			} else {
				err = c.Call(nombreServidor+metodoCalifiacionesMaterias, "", &respuesta)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(respuesta)
					fmt.Println("Promedio general = " + fmt.Sprintf("%f", promedio))
				}
			}
		case 4:
			fmt.Println("\t---Obtener promedio de una materia---")
			nombreMateria := dameNombreMateria()
			err = c.Call(nombreServidor+metodoPromedioMateria, nombreMateria, &promedio)
			if err != nil {
				fmt.Println(err)
			} else {
				err = c.Call(nombreServidor+metodoCalificacionesAlumnosMateria, nombreMateria, &respuesta)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(respuesta)
					fmt.Println("Promedio de la materia = " + fmt.Sprintf("%f", promedio))
				}
			}
		case 0:
			c.Close()
			continuar = false
			fmt.Println("Finalizar servidor")
		default:
			fmt.Println("\nOpción no válida")
		}
		if continuar {
			var enter string
			fmt.Println("\nPresione enter para continuar . . .")
			fmt.Scanln(&enter)
		}
	}
}

func agregarCalificacionAlumno() Mensaje {
	var calificacion float64
	fmt.Println("\n\t---Agregar la calificación de una materia para un alumno---")
	fmt.Print("\nIngrese nombre del alumno: ")
	in := bufio.NewReader(os.Stdin)
	nombreAlumno, _ := in.ReadString('\n')
	nombreAlumno = nombreAlumno[0 : len(nombreAlumno)-2]
	fmt.Print("Ingrese nombre de la materia: ")
	nombreMateria, _ := in.ReadString('\n')
	nombreMateria = nombreMateria[0 : len(nombreMateria)-2]
	fmt.Print("Ingrese calificación: ")
	fmt.Scanln(&calificacion)
	return Mensaje{nombreMateria, nombreAlumno, calificacion}
}

func dameNombreAlumno() string {
	fmt.Print("\nIngrese nombre del alumno: ")
	in := bufio.NewReader(os.Stdin)
	nombreAlumno, _ := in.ReadString('\n')
	nombreAlumno = nombreAlumno[0 : len(nombreAlumno)-2]
	return nombreAlumno
}

func dameNombreMateria() string {
	fmt.Print("\nIngrese nombre de la materia: ")
	in := bufio.NewReader(os.Stdin)
	nombreMateria, _ := in.ReadString('\n')
	nombreMateria = nombreMateria[0 : len(nombreMateria)-2]
	return nombreMateria
}