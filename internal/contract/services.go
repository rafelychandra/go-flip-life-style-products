package contract

import "go-flip-life-style-products/internal/services"

type Service struct {
	Balance     services.Balance
	Statements  services.Statements
	Transaction services.Transaction
}

func newService(c *Contract) (*Service, error) {
	s := &Service{}
	s.Balance = services.NewBalance(c.Cfg, c.Repositories.Transaction)
	s.Statements = services.NewStatements(c.Cfg, c.File, c.Queue)
	s.Transaction = services.NewTransaction(c.Cfg, c.Repositories.Transaction, c.Event)
	return s, nil
}
