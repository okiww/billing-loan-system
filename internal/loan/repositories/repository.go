package repositories

import "sync"

var (
	repoLoan         LoanRepositoryInterface
	repoLoanBill     LoanBillRepositoryInterface
	repoLoanLock     sync.Once
	repoLoanBillLock sync.Once
)
