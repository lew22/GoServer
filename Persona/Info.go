package Persona

type Info struct {
	Nombre   string
	Apellido string
	Opcion   string
	Metodo   string
}

func NuevaInfo(nombre string, apellido string, opcion string, metodo string) *Info {
	return &Info{nombre, apellido, opcion, metodo}
}
