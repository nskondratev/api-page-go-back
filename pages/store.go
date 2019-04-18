package pages

type Store interface {
	GetById(uint64) (*Page, error)
	List(offset, limit int, sort string, descending bool, query string) ([]*PageList, int, error)
	Update(*Page) error
	Delete(*Page) error
	Create(*Page) error
}
