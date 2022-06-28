package Persona

type Info struct {
	Nombre   string
	Apellido string
	Opcion   string
}

func NuevaInfo(nombre string, apellido string, opcion string) *Info {
	return &Info{nombre, apellido, opcion}
}
