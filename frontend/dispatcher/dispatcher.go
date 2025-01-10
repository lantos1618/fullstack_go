package dispatcher

// ID is a unique identifier for registered callbacks
type ID int

var (
	idCounter ID
	callbacks = make(map[ID]func(action interface{}))
)

// Dispatch dispatches an action to all registered callbacks
func Dispatch(action interface{}) {
	for _, c := range callbacks {
		c(action)
	}
}

// Register registers a callback to handle dispatched actions
func Register(callback func(action interface{})) ID {
	idCounter++
	id := idCounter
	callbacks[id] = callback
	return id
}

// Unregister removes a previously registered callback
func Unregister(id ID) {
	delete(callbacks, id)
}
