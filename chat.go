package ollama

// Chat stores the messages sent from the user and received from the assistant.
type Chat struct {
	ID       string
	Messages []Message
}

// AddMessage adds a new message to the end of the chat.
//
// Parameters:
//   - m: The message to add.
func (c *Chat) AddMessage(m Message) {
	c.Messages = append(c.Messages, m)
}

// AddMessageTo adds a new message at the specified index.
//
// Parameters:
//   - index: The index at which to add the new message.
//   - m: The message to add.
func (c *Chat) AddMessageTo(index int, m Message) {
	c.Messages = append(c.Messages[:index], append([]Message{m}, c.Messages[index:]...)...)
}

// DeleteMessage deletes a message at the specified index.
//
// Parameters:
//   - index: The index of the message to delete.
func (c *Chat) DeleteMessage(index int) {
	c.Messages = append(c.Messages[:index], c.Messages[index+1:]...)
}

// DeleteAllMessages deletes all messages in the chat.
func (c *Chat) DeleteAllMessages() {
	c.Messages = make([]Message, 0)
}
