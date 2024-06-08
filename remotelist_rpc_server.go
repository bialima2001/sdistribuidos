package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"
	"ppgti/remotelist/pkg"
)

// RemoteList define o componente que gerencia listas de valores inteiros de forma remota.
type RemoteList struct {
	mu    sync.Mutex
	lists map[string][]int
}

// Append adiciona um valor ao final da lista com o identificador listID.
func (rl *RemoteList) Append(args *ListArgs, reply *bool) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	list, ok := rl.lists[args.ListID]
	if !ok {
		list = []int{}
	}

	list = append(list, args.Value)
	rl.lists[args.ListID] = list

	*reply = true

	// Salva o estado atual das listas em um arquivo
	err := rl.saveState()
	if err != nil {
		return err
	}

	return nil
}

// Get retorna o valor na posição index da lista com o identificador listID.
func (rl *RemoteList) Get(args *ListArgs, reply *int) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	list, ok := rl.lists[args.ListID]
	if !ok || args.Index >= len(list) {
		return fmt.Errorf("index out of range or list not found: %v", args.ListID)
	}

	*reply = list[args.Index]
	return nil
}

// Remove remove o último elemento da lista com o identificador listID e retorna o valor removido.
func (rl *RemoteList) Remove(args *ListArgs, reply *int) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	list, ok := rl.lists[args.ListID]
	if !ok || len(list) == 0 {
		return fmt.Errorf("list is empty or list not found: %v", args.ListID)
	}

	lastIndex := len(list) - 1
	lastValue := list[lastIndex]
	rl.lists[args.ListID] = list[:lastIndex]

	*reply = lastValue

	// Salva o estado atual das listas em um arquivo
	err := rl.saveState()
	if err != nil {
		return err
	}

	return nil
}

// Size retorna o tamanho da lista com o identificador listID.
func (rl *RemoteList) Size(args *ListArgs, reply *int) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	list, ok := rl.lists[args.ListID]
	if !ok {
		return fmt.Errorf("list not found: %v", args.ListID)
	}

	*reply = len(list)
	return nil
}

// ListArgs contém os argumentos dos métodos RPC.
type ListArgs struct {
	ListID string // Identificador da lista
	Value  int    // Valor a ser adicionado à lista
	Index  int    // Índice do elemento a ser recuperado da lista
}

// Método para salvar o estado atual das listas em um arquivo
func (rl *RemoteList) saveState() error {
	file, err := os.Create("lists.gob")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(rl.lists)
	if err != nil {
		return err
	}

	return nil
}

// Método para carregar o estado das listas a partir de um arquivo
func (rl *RemoteList) loadState() error {
	file, err := os.Open("lists.gob")
	if err != nil {
		if os.IsNotExist(err) {
			// O arquivo não existe, então não há estado para carregar
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&rl.lists)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Cria uma nova instância de RemoteList
	remoteList := &RemoteList{
		lists: make(map[string][]int),
	}

	// Carrega o estado das listas a partir do arquivo (se existir)
	err := remoteList.loadState()
	if err != nil {
		fmt.Println("Erro ao carregar o estado das listas:", err)
	}

	// Registra a instância RemoteList como um serviço RPC
	rpc.Register(remoteList)

	// Configura o servidor RPC para ouvir na porta 1234
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Erro ao ouvir na porta 1234:", err)
		return
	}
	defer l.Close()

	fmt.Println("RemoteList server está rodando na porta 1234...")
	for {
		conn, _ := l.Accept()
		go rpc.ServeConn(conn)
	}
}
