package page

// Constants for the different message types, they correspond to bootstrap 4
// alert classes
const (
	MessageInfo    = "info"
	MessageSuccess = "success"
	MessageWarning = "warning"
	MessageError   = "danger"
)

// Message is
type Message struct {
	Type        string
	Message     string
	Dismissable bool
}

// NewMessage creates a message struct which can be used to show messages to the
// end user.
func NewMessage(typ, mess string, dismis bool) *Message {
	m := Message{
		typ, mess, dismis,
	}

	return &m
}
