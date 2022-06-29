package Estructura

import (
	"fmt"

	"tf.com/events/Persona"
)

//se usa ../ cuando se llama de una carpeta
// ./ cuando se llama fuera de una carpeta
type Nodo struct {
	siguiente *Nodo
	info      *Persona.Info
}

type Lista struct {
	primero *Nodo
	ultimo  *Nodo
	Cont    int
}

func NuevaLista() *Lista {
	return &Lista{nil, nil, 0}
}

func CrearNodo(info *Persona.Info) *Nodo {
	//referencia de donde esta la informacion
	return &Nodo{siguiente: nil, info: info}
}

func Insertar(info *Persona.Info, lista *Lista) {
	var nuevo *Nodo = CrearNodo(info)

	//lista vacia
	if lista.primero == nil {
		lista.primero = nuevo
		lista.ultimo = nuevo
		lista.Cont += 1
	} else {
		lista.ultimo.siguiente = nuevo
		lista.ultimo = lista.ultimo.siguiente
		lista.ultimo.siguiente = lista.primero
		lista.Cont += 1
	}

}

func Imprimir(lista *Lista) {
	aux := lista.primero

	for aux != nil {
		fmt.Println("{")
		fmt.Println("Nombre: ", aux.info.Nombre)
		fmt.Println("Apellido: ", aux.info.Apellido)
		fmt.Println("Opcion: ", aux.info.Opcion)
		fmt.Println("Metodo: ", aux.info.Metodo)
		fmt.Println("}")
		aux = aux.siguiente
	}

}
