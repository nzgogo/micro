package context

type context struct {
	pool []*Conversation
}

type Conversation struct {
	done chan int
}
