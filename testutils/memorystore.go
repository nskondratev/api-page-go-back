package testutils

import (
	"github.com/nskondratev/api-page-go-back/util"
)

func NewArrayNullStringFromStrings(in []string) ([]util.NullString, error) {
	res := make([]util.NullString, len(in), len(in))
	for i, str := range in {
		ns, err := util.NewNullStringFromString(str)

		if err != nil {
			return nil, err
		}

		res[i] = ns
	}
	return res, nil
}

type MemoryCreateTestCase struct {
	ItemToCreate interface{}
	TotalRows    int
	LastItemID   uint64
}

type MemoryDeleteTestCase struct {
	ItemToDelete interface{}
	TotalRows    int
}

type MemoryGetByIdTestCase struct {
	ID           uint64
	ExpectedItem interface{}
}

type MemoryUpdateTestCase struct {
	ItemToUpdate interface{}
}
