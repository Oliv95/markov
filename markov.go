// Package markov provides a very basic interface for creating markov chains
package markov

import "fmt"
import "errors"
import "math/rand"
import "time"
import "sync"

// State represents a state in a markov chain.
// The Data field needs be comparable, it will be used as a map key
type State struct {
	Data interface{}
}

// Graph represents the markov chain
type Graph struct {
	table map[State][]State
	lock  *sync.RWMutex
}

// CreateEmptyGraph returns a empty graph
func CreateEmptyGraph() Graph {
	graph := Graph{
		map[State][]State{},
		&sync.RWMutex{},
	}
	return graph

}

// AddTransition adds the transition
// 	from -> to
// into the markov chain or if the transition already exists, increases it's weight
func AddTransition(graph *Graph, from State, to State) {
	// TODO change this to only store ratio, no need to store duplicates
	// TODO make thread safe
	graph.lock.Lock()
	defer graph.lock.Unlock()
	graph.table[from] = append(graph.table[from], to)
}

// RandomState returns a random state from the table
// or if the table is empty or nil errors
func RandomState(graph *Graph) (State, error) {
	emptyOrNil := graph == nil || len(graph.table) == 0
	if emptyOrNil {
		return State{}, fmt.Errorf("Graph cannot be empty or nil: %v", graph.table)
	}
	rand.Seed(time.Now().UTC().UnixNano())

	graph.lock.RLock()
	defer graph.lock.RUnlock()

	keys := make([]State, 0, len(graph.table))
	for k := range graph.table {
		keys = append(keys, k)
	}
	nbrKeys := len(keys)
	randomIndex := rand.Intn(nbrKeys)
	randomState := keys[randomIndex]
	return randomState, nil
}

// Transition returns the new state after one time step from the given state
//
// Errors if the table is nil, current does not exist in the table or does not have any transitions defined
func Transition(graph *Graph, current State) (*State, error) {
	if graph == nil {
		return nil, errors.New("Graph cannot be nil")
	}

	graph.lock.RLock()
	defer graph.lock.RUnlock()

	transitions, ok := (graph.table)[current]
	if !ok {
		return nil, fmt.Errorf("State %v not in table %v", current, graph.table)
	}
	nbrTransitions := len(transitions)
	// Checks if there is no possible state to move to
	if nbrTransitions == 0 {
		return nil, fmt.Errorf("No transitions for state %v", current)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	randomIndex := rand.Intn(nbrTransitions)
	return &transitions[randomIndex], nil
}
