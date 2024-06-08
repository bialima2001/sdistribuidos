package main

import (
	"fmt"
	"log"
	"net/rpc"
	"ppgti/remotelist/pkg"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Erro ao conectar:", err)
	}
	defer client.Close()

	var reply bool
	args := &remotelist.ListArgs{
		ListID: "mylist",
		Value:  42,
	}
	err = client.Call("RemoteList.Append", args, &reply)
	if err != nil {
		log.Fatal("Erro ao chamar Append:", err)
	}
	fmt.Println("Append result:", reply)

	var getValue int
	getArgs := &remotelist.ListArgs{
		ListID: "mylist",
		Index:  0,
	}
	err = client.Call("RemoteList.Get", getArgs, &getValue)
	if err != nil {
		log.Fatal("Erro ao chamar Get:", err)
	}
	fmt.Println("Get result:", getValue)

}
