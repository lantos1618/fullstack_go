package actions

// SetUsername is an action that sets the current user's username
type SetUsername struct {
	Username string
}

// AddMessage is an action that adds a new chat message
type AddMessage struct {
	Text string
	From string
}

// SetTyping is an action that sets a user's typing status
type SetTyping struct {
	Username string
	IsTyping bool
}
