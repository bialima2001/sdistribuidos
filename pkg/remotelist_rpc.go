package remotelist

import (
	"errors"
	"fmt"
	"sync"
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
	return nil
}

// Get retorna o valor na posição index da lista com o identificador listID.
func (rl *RemoteList) Get(args *ListArgs, reply *int) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	list, ok := rl.lists[args.ListID]
	if !ok || args.Index >= len(list) {
		return errors.New("index out of range")
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
		return errors.New("list is empty")
	}

	lastIndex := len(list) - 1
	lastValue := list[lastIndex]
	rl.lists[args.ListID] = list[:lastIndex]

	*reply = lastValue
	return nil
}

// Size retorna o tamanho da lista com o identificador listID.
func (rl *RemoteList) Size(args *ListArgs, reply *int) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	list, ok := rl.lists[args.ListID]
	if !ok {
		return errors.New("list not found")
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
