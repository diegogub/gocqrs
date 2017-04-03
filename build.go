package gocqrs

type ViewBuilder struct {
	stream        string
	Name          string `json:"name"`
	CurrentStatus uint64 `json:"currentStatus"`
	Failed        bool   `json:"failed"`
	LastError     string `json:"lastError"`
	Store         EventStore
	Handler       Event
}

type EventApplier interface {
	Apply(e Event) (uint64, error)
}

func NewViewBuilder(name string) *ViewBuilder {
	var vb ViewBuilder

	return &vb
}
