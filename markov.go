// Very basic package for creating markov chains
package markov

import "fmt"
import "errors"
import "math/rand"
import "time"

// Represents a state in a markov chain
type State struct {
	Data interface{}
}

// A table that stores all the transitions
// Each state is assosicated with the all transitions from that state
type TransitionTable = map[State][]State

// Returns a empty table
func CreateEmptyTable() TransitionTable {
	return make(TransitionTable)

}

// Adds the transition
// 	from -> to
// into the markov chain or if the transition already exists, increases it's weight
func AddTransition(table *TransitionTable, from State, to State) {
	// TODO change this to only store ratio, no need to store duplicates
	// TODO make thread safe
	(*table)[from] = append((*table)[from], to)
}

// Returns a random state from the table
// or if the table is empty or nil errors
func RandomState(table *TransitionTable) (State, error) {
	emptyOrNil := table == nil || len(*table) == 0
	if emptyOrNil {
		return State{}, errors.New("Table cannot be empty or nil")
	}
	rand.Seed(time.Now().UTC().UnixNano())
	keys := make([]State, 0, len(*table))
	for k, _ := range *table {
		keys = append(keys, k)
	}
	nbrKeys := len(keys)
	randomIndex := rand.Intn(nbrKeys)
	randomState := keys[randomIndex]
	return randomState, nil
}

// Returns the new state after one time step from the given state
//
// Errors if the table is nil, current does not exist in the table or does not have any transitions defined
func Transition(table *TransitionTable, current State) (*State, error) {
	if table == nil {
		return nil, errors.New("Table cannot be nil")
	}
	transitions, ok := (*table)[current]
	if !ok {
		return nil, errors.New(fmt.Sprint("State %v not in table %v", current, *table))
	}
	nbrTransitions := len(transitions)
	// Checks if there is no possible state to move to
	if nbrTransitions == 0 {
		return nil, errors.New(fmt.Sprint("No transitions for state %v", current))
	}
	rand.Seed(time.Now().UTC().UnixNano())
	randomIndex := rand.Intn(nbrTransitions)
	return &transitions[randomIndex], nil
}
