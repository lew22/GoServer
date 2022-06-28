# GoServer

terminal:

Para usar los packages ,activamos los modulos

go env -w GO111MODULE=on

////////////////////////////////////////////////
go build node.go

node localhost:port

ó

node (fuente):port (destino):port

ejemplo:

node localhost:8001 localhost:8000

node 127.0.0.1:5500 localhost:8000

ó

node localhost:8003 localhost:8001 "text"

ejemplo

node localhost:8003 localhost:8001 goo

node localhost:8003 localhost:8001 agrawalla
