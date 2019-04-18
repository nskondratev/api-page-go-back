package events

type Store interface {
	GetById(uint64) (*Event, error)
	List(offset, limit int, sort string, descending bool, eType string, query string) ([]*EventList, int, error)
	Create(*Event) error
	Update(*Event) error
	Delete(*Event) error
}
