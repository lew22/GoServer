package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"tf.com/events/Estructura"
	"tf.com/events/Persona"
)

type Frame struct {
	Cmd    string   `json:"cmd"`
	Sender string   `json:"sender"`
	Data   []string `json:"data"`
}

type Info struct {
	nextNode string
	nextNum  int
	imFirst  bool
	cont     int
}

type ConsInfo struct {
	contA int
	contB int
}

var (
	host         string
	myNum        int
	chRemotes    chan []string
	chInfo       chan Info
	chCons       chan ConsInfo
	participants int
	readyToStart chan bool
	gotnums      bool
)

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) == 1 {
		log.Println("Hostname not given")
	} else {
		chRemotes = make(chan []string, 1)
		chInfo = make(chan Info, 1)
		chCons = make(chan ConsInfo, 1)
		readyToStart = make(chan bool, 1)

		host = os.Args[1]
		chRemotes <- []string{}
		if len(os.Args) >= 3 {
			connectToNode(os.Args[2])
		}
		if len(os.Args) == 4 {
			switch os.Args[3] {
			case "agrawalla":
				go startAgrawalla()
			case "consensus":
				go startConsensus()
			}
		}
		server()
		//GenerarNodos()

	}

}

//generamos los nodos
// func GenerarNodos() {
// 	var info1 *Persona.Info = Persona.NuevaInfo("richard", "garcia", "A")
// 	var lista *Estructura.Lista = Estructura.NuevaLista()

// 	Estructura.Insertar(info1, lista)
// 	Estructura.Imprimir(lista)
// }

//decodificar json
func Decoficacion(data string) {
	var info1 Persona.Info

	dataB := []byte(data)
	err := json.Unmarshal(dataB, &info1)

	if err != nil {
		fmt.Println("Error decodificando:", err)
	} else {
		// fmt.Println("Nombre: ", info1.Nombre)
		// fmt.Println("Apellido: ", info1.Apellido)
		// fmt.Println("Opcion: ", info1.Opcion)
		fmt.Println("Decodificacion Exitosa")
		CrearNodo(info1.Nombre, info1.Apellido, info1.Opcion)
	}

}

func CrearNodo(nombre string, apellido string, opcion string) {
	var Ninfo *Persona.Info = Persona.NuevaInfo(nombre, apellido, opcion)
	var lista *Estructura.Lista = Estructura.NuevaLista()

	Estructura.Insertar(Ninfo, lista)
	fmt.Println("Se insertaron los nodos y esta es la estructura")
	Estructura.Imprimir(lista)
}

func startAgrawalla() {
	go func() {
		time.Sleep(5 * time.Second)
		gotnums = false
		remotes := <-chRemotes
		chRemotes <- remotes
		for _, remote := range remotes {
			send(remote, Frame{"agrawalla", host, []string{}}, nil)
		}
		handleAgrawalla()
	}()
}
func startConsensus() {
	remotes := <-chRemotes
	for _, remote := range remotes {
		log.Printf("%s: sending consensus to %s\n", host, remote)
		send(remote, Frame{"consensus", host, []string{}}, nil)
	}
	chRemotes <- remotes
	handleConsensus()
}

func timeOutChecker(seconds int) {
	count := seconds * 10
	for i := 0; i < count; i++ {
		time.Sleep(100 * time.Millisecond)
		if gotnums {
			break
		}
	}
	if !gotnums {
		log.Printf("%s: timeout!!", host)
	}
}

func connectToNode(remote string) {
	remotes := <-chRemotes
	remotes = append(remotes, remote)
	chRemotes <- remotes
	if !send(remote, Frame{"hello", host, []string{}}, func(cn net.Conn) {
		dec := json.NewDecoder(cn)
		var frame Frame
		dec.Decode(&frame)
		remotes := <-chRemotes
		remotes = append(remotes, frame.Data...)
		chRemotes <- remotes
		log.Printf("%s: friends: %s\n", host, remotes)
	}) {
		log.Printf("%s: unable to connect to %s\n", host, remote)
	}
}

func send(remote string, frame Frame, callback func(net.Conn)) bool {
	if cn, err := net.Dial("tcp", remote); err == nil {
		defer cn.Close()
		enc := json.NewEncoder(cn)
		enc.Encode(frame)
		if callback != nil {
			callback(cn)
		}
		return true
	} else {
		log.Printf("%s: can't connect to %s\n", host, remote)
		idx := -1
		remotes := <-chRemotes
		for i, rem := range remotes {
			if remote == rem {
				idx = i
				break
			}
		}
		if idx >= 0 {
			remotes[idx] = remotes[len(remotes)-1]
			remotes = remotes[:len(remotes)-1]
		}
		chRemotes <- remotes
		return false
	}
}

func server() {
	if ln, err := net.Listen("tcp", host); err == nil {
		defer ln.Close()
		log.Printf("Listening on %s\n", host)
		for {
			cn, err := ln.Accept()
			if err == nil {
				go fauxDispatcher(cn)
			} else {
				log.Printf("%s: cant accept connection.\n", host)
			}
			//procesar una conexion
			go ProcesarCliente(cn)
		}
	} else {
		log.Printf("Can't listen on %s\n", host)
	}
}

//procesamos al cliente cuando se establece la conexion
func ProcesarCliente(cn net.Conn) {
	buffer := make([]byte, 1024)
	pc, err := cn.Read(buffer)

	if err != nil {
		fmt.Println("Oops ocurrio un error:", err.Error())
	}
	fmt.Println("Recibimos conexion: ", string(buffer[:pc]))
	_, err = cn.Write([]byte("Hola, te estanmos respondiendo"))
	//+ string(buffer[:pc])

	if err != nil {
		fmt.Println("Oops ocurrio un error:", err.Error())
	} else {
		Decoficacion(string(buffer[:pc]))
	}
	cn.Close()

}

func fauxDispatcher(cn net.Conn) {
	defer cn.Close()
	dec := json.NewDecoder(cn)
	frame := &Frame{}
	dec.Decode(frame)
	switch frame.Cmd {
	case "hello":
		handleHello(cn, frame)
	case "add":
		handleAdd(frame)
	case "agrawalla":
		handleAgrawalla()
	case "num":
		handleNum(frame)
	case "start":
		handleStart()
	case "consensus":
		handleConsensus()
	case "vote":
		handleVote(frame)
		/*case "potato":
		handlePotato(frame)*/
	}
}

func handleHello(cn net.Conn, frame *Frame) {
	enc := json.NewEncoder(cn)
	remotes := <-chRemotes
	enc.Encode(Frame{"<response>", host, remotes})
	notification := Frame{"add", host, []string{frame.Sender}}
	for _, remote := range remotes {
		send(remote, notification, nil)
	}
	remotes = append(remotes, frame.Sender)
	log.Printf("%s: friends: %s\n", host, remotes)
	chRemotes <- remotes
}
func handleAdd(frame *Frame) {
	remotes := <-chRemotes
	remotes = append(remotes, frame.Data...)
	log.Printf("%s: friends: %s\n", host, remotes)
	chRemotes <- remotes
}
func handleAgrawalla() {
	myNum = rand.Intn(1000000000)
	go timeOutChecker(10)
	log.Printf("%s: my number is %d\n", host, myNum)
	msg := Frame{"num", host, []string{strconv.Itoa(myNum)}}
	remotes := <-chRemotes
	chRemotes <- remotes
	for _, remote := range remotes {
		send(remote, msg, nil)
	}
	chInfo <- Info{"", 1000000001, true, 0}
}
func handleNum(frame *Frame) {
	if num, err := strconv.Atoi(frame.Data[0]); err == nil {
		info := <-chInfo
		//log.Printf("from %v\n", frame)
		if num > myNum {
			if num < info.nextNum {
				info.nextNum = num
				info.nextNode = frame.Sender
			}
		} else {
			info.imFirst = false
		}
		info.cont++
		chInfo <- info
		remotes := <-chRemotes
		chRemotes <- remotes
		if info.cont == len(remotes) {
			if info.imFirst {
				log.Printf("%s: I'm first!\n", host)
				criticalSection()
			} else {
				readyToStart <- true
			}
			gotnums = true
		}
	} else {
		log.Printf("%s: can't convert %v\n", host, frame)
	}
}
func handleStart() {
	<-readyToStart
	criticalSection()
}
func handleConsensus() {
	time.Sleep(3 * time.Second)
	fmt.Print("A o B, elige una: ")
	var op string
	fmt.Scanf("%s\n", &op)
	info := ConsInfo{0, 0}
	if op == "A" {
		info.contA++
	} else {
		info.contB++
	}
	chCons <- info
	remotes := <-chRemotes
	participants = len(remotes) + 1
	for _, remote := range remotes {
		log.Printf("%s: sending %s to %s\n", host, op, remote)
		send(remote, Frame{"vote", host, []string{op}}, nil)
	}
	chRemotes <- remotes
}
func handleVote(frame *Frame) {
	vote := frame.Data[0]
	info := <-chCons
	if vote == "A" {
		info.contA++
	} else {
		info.contB++
	}
	chCons <- info
	log.Printf("%s: got %v\n", host, frame)
	if info.contA+info.contB == participants {
		if info.contA > info.contB {
			log.Printf("%s go A\n", host)
		} else {
			log.Printf("%s go B\n", host)
		}
	}
}

/*
func handlePotato(frame *Frame) {
	if num, err := strconv.Atoi(frame.Data[0]); err == nil {
		log.Printf("%s: recibí %d\n", host, num)
		if num == 0 {
			log.Printf("%s: perdí\n", host)
		} else {
			for len(remotes) > 0 {
				remote := remotes[rand.Intn(len(remotes))]
				data := []string{strconv.Itoa(num - 1)}
				time.Sleep(100 * time.Millisecond)
				if send(remote, Frame{"potato", host, data}, nil) {
					break
				}
			}
		}
	} else {
		log.Printf("%s: can't convert %s to number\n", host, frame.Data)
	}
}
*/
func criticalSection() {
	log.Printf("%s: my time has come!\n", host)
	info := <-chInfo
	if info.nextNode != "" {
		log.Printf("%s: letting %s start\n", host, info.nextNode)
		send(info.nextNode, Frame{"start", host, []string{}}, nil)
	} else {
		log.Printf("%s: I was the last one :(\n", host)
	}
}

/*
func potatoGenerator() {
	for {
		time.Sleep(5 * time.Second)
		for len(remotes) > 0 {
			remote := remotes[rand.Intn(len(remotes))]
			data := []string{strconv.Itoa(rand.Intn(20) + 10)}
			if send(remote, Frame{"potato", host, data}, nil) {
				break
			}
		}
	}
}
*/
