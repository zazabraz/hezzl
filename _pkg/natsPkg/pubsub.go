package natsPkg

type Queue[MessageType Message] struct {
	PubQueue[MessageType]
	SubQueue[MessageType]
}

func NewQueue[MessageType Message](pubQueue PubQueue[MessageType], subQueue SubQueue[MessageType]) *Queue[MessageType] {
	return &Queue[MessageType]{PubQueue: pubQueue, SubQueue: subQueue}
}
