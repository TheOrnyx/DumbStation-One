package gpu

import "github.com/TheOrnyx/psx-go/log"

// Buffer holding multi-word GP0 command parameters
//
// TODO - idk if this needed to be a whole type i was just following the guide but idk
type CommandBuffer struct {
	buffer [12]uint32 // the buffer - longest command is GP0(3Eh) which takes 12 params
	length uint8      // the Number of words queued in the buffer
}

// newBuffer CommandBuffer constructor
func newBuffer() CommandBuffer {
	return CommandBuffer{buffer: [12]uint32{}, length: 0}
}

// clear clear the command buffer
func (c *CommandBuffer) clear()  {
	c.length = 0
}

// pushWord push word to the buffer
func (c *CommandBuffer) pushWord(word uint32)  {
	c.buffer[c.length] = word
	c.length += 1
}

// at return the command at location index in buffer
func (c *CommandBuffer) at(index uint8) uint32 {
	if index >= c.length {
		log.Panicf("Command buffer index out of range %v with length %v", index, c.length)
	}

	return c.buffer[index]
}
