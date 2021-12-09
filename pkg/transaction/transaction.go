package transaction

type Transaction struct {
	parent  *Transaction
	actions map[string]interface{}
}

func NewTransaction(previousTransaction *Transaction) Transaction {
	return Transaction{
		parent:  previousTransaction,
		actions: make(map[string]interface{}),
	}
}
