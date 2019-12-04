package main

import (
	"log"
	"markV3/mface"
	"markV3/mnet"
	"time"
)

func entrance() error {
	log.Println("This is entrance func")
	return nil
}

func routes() map[string]func(mface.MMessage, mface.MMessage) error {
	routes := map[string]func(mface.MMessage, mface.MMessage) error {}

	routes["0000000000"] = firstHandleFunc
	routes["0000000001"] = secondHandleFunc

	return routes
}

func main() {
	server , err := mnet.NewServer()
	if err != nil {
		panic(err)
	}

	_ = server.AddRoutes(routes())

	server.RunEntranceFunc(entrance)

	log.Println(server.Start())

	time.Sleep(time.Second * time.Duration(3))

	log.Println(server.Stop())

}

func firstHandleFunc(request mface.MMessage, response mface.MMessage) error {
	log.Printf("This is \"0000000000\" route handleFunc : \nreuqest : %+v \nresponse : %+v" , request , response)
	return nil
}

func secondHandleFunc(request mface.MMessage, response mface.MMessage) error {
	log.Printf("This is \"0000000001\" route handleFunc : \nreuqest : %+v \nresponse : %+v" , request , response)
	return nil
}