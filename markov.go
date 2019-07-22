// Package markov provides a very basic interface for creating markov chains
package markov

import "fmt"
import "errors"
import "math/rand"
import "time"

// State represents a state in a markov chain.
// The Data field needs be comparable, it will be used as a map key
type State struct {
	Data interface{}
}

// TransitionTable stores all the transitions in the markov chain
// Each state is assosicated with the all transitions from that state
type TransitionTable = map[State][]State

// CreateEmptyTable returns a empty transitionTable
func CreateEmptyTable() TransitionTable {
	return make(TransitionTable)

}

// AddTransition adds the transition
// 	from -> to
// into the markov chain or if the transition already exists, increases it's weight
func AddTransition(table *TransitionTable, from State, to State) {
	// TODO change this to only store ratio, no need to store duplicates
	// TODO make thread safe
	(*table)[from] = append((*table)[from], to)
}

// RandomState returns a random state from the table
// or if the table is empty or nil errors
func RandomState(table *TransitionTable) (State, error) {
	emptyOrNil := table == nil || len(*table) == 0
	if emptyOrNil {
		return State{}, fmt.Errorf("Table cannot be empty or nil: %v", table)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	keys := make([]State, 0, len(*table))
	for k := range *table {
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
func Transition(table *TransitionTable, current State) (*State, error) {
	if table == nil {
		return nil, errors.New("Table cannot be nil")
	}
	transitions, ok := (*table)[current]
	if !ok {
		return nil, fmt.Errorf("State %v not in table %v", current, *table)
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
