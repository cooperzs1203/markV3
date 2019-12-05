package main

import (
	"log"
	"markV3/test"
)

//
//import (
//	"log"
//	"markV3/mface"
//	"markV3/mnet"
//	"net"
//	"time"
//)
//
//var dp mface.MDataProtocol
//
//func init() {
//	dp = mnet.NewDataProtocol("HEAD", 10, 4, "", 100)
//}
//
//func entrance() error {
//	log.Println("This is entrance func")
//	return nil
//}
//
//func routes() map[string]mface.RouteHandleFunc {
//	routes := map[string]mface.RouteHandleFunc{}
//
//	routes["0000000000"] = firstHandleFunc
//	routes["0000000001"] = secondHandleFunc
//
//	return routes
//}
//
//func main() {
//	server, err := mnet.NewServer()
//	if err != nil {
//		panic(err)
//	}
//
//	if err := server.AddRoutes(routes()); err != nil {
//		log.Println("server.AddRoutes error", err)
//		return
//	}
//
//	server.AddRequestHook(customRequestHook)
//	server.AddResponseHook(customResponseHook)
//
//	server.RunEntranceFunc(entrance)
//
//	err = server.Start()
//	if err != nil {
//		panic(err)
//	}
//
//	go clientTest()
//
//	for {
//		time.Sleep(time.Second * time.Duration(10))
//	}
//}
//
//func customRequestHook(request mface.MMessage) bool {
//	log.Printf("This is request hook func : \nreuqest : %+v ", request)
//
//	return true
//}
//
//func customResponseHook(response mface.MMessage) bool {
//	log.Printf("This is response hook func : \nreuqest : %+v ", response)
//
//	return true
//}
//
//func firstHandleFunc(request mface.MMessage) mface.MMessage {
//	log.Printf("This is \"0000000000\" route handleFunc : \nreuqest : %+v", request)
//
//	//var response mface.MMessage
//	//dp.Unmarshal()
//
//	return nil
//}
//
//func secondHandleFunc(request mface.MMessage) mface.MMessage {
//	log.Printf("This is \"0000000001\" route handleFunc : \nreuqest : %+v", request)
//	return nil
//}
//
//func clientTest() {
//	conn, err := net.Dial("tcp", "0.0.0.0:8888")
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	go func() {
//
//		for {
//			time.Sleep(time.Second * time.Duration(5))
//			data := dp.Marshal("0000000000", []byte("From Cooper --MarkV3 Test."))
//
//			if _, err := conn.Write(data); err != nil {
//				log.Println("write data err : ", err)
//				continue
//			}
//		}
//	}()
//
//	for {
//		buffer := make([]byte, 1024)
//		cnt, err := conn.Read(buffer)
//		if err != nil {
//			log.Println("read error : ", err)
//			continue
//		}
//
//		log.Println("Get Response : ", string(buffer[18:cnt]))
//	}
//}

func main() {
	dp, err := test.NewDataProtocol("HEAD", 10 , 4)
	if err != nil {
		log.Println(err)
		return
	}

	//data := []byte("fuck you from cooper")
	//cmData , err := dp.EnPack([]byte("0000000000") , []byte("KK"), test.IntToByte(uint32(len(data))) , data)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//log.Println(cmData)
	//log.Println(string(cmData))

	prefix := []byte{1, 2, 3, 4, 5}
	suffix := []byte{6, 7, 8, 9, 0}
	//data := []byte{72,69,65,68,48,48,48,48,48,48,48,48,48,48,75,75,0,0,0,20,102,117,99,107,32,121,111,117,32,102,114,111,109,32,99,111,111,112,101,114}
	data := []byte{72, 69, 65, 68, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 0, 0, 0, 20, 102, 117, 99, 107, 32, 121, 111, 117, 32, 102, 114, 111, 109, 32, 99, 111, 111, 112, 101, 114}

	data = append(prefix, data...)
	data = append(data, suffix...)

	go func() {
		for {
			msg, ok := <-dp.CompletedDataChan()
			if !ok {
				break
			}
			log.Println(msg.Data())
			log.Println(string(msg.Data()))
			log.Println(msg.DataLength())
			for i := uint32(0); i < msg.ValuesLength(); i++ {
				log.Println(msg.Value(int(i)))
			}

			log.Println(msg.Marshal())
		}
	}()

	dp.DePack(data)

	for {

	}



}
