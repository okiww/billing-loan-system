package repositories

import "sync"

var (
	repoLoan     LoanRepositoryInterface
	repoLoanBill LoanBillRepositoryInterface
	repoLock     sync.Once
)
