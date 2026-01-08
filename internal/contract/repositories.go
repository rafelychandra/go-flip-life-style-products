package contract

import "go-flip-life-style-products/internal/repositories"

type Repositories struct {
	Transaction repositories.Transaction
}

func newRepository(c *Contract) (*Repositories, error) {
	r := &Repositories{}

	r.Transaction = repositories.NewTransaction()
	return r, nil
}
