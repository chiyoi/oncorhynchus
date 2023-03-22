package sakana

type Handler interface {
	Serve(arguments []string)
}

type HandlerFunc func(arguments []string)

func (hf HandlerFunc) Serve(arguments []string) { hf(arguments) }
