serve-http:
	@go run main.go http

lint:
	@golangci-lint run -E gofmt

format:
	@$(MAKE) fmt
	@$(MAKE) imports

fmt:
	@echo "Formatting code style..."
	gofmt -w -s cmd/.. \
		configs/.. \
		internal/..
	@echo "[DONE] Formatting code style..."

# Get the root directory of the project dynamically using git
REPO_PATH := $(shell git rev-parse --show-toplevel)

imports:
	@echo "Formatting imports..."
	# Use the dynamically determined repository path
	goimports -w -local $(REPO_PATH)/billing-loan-system \
		cmd/.. \
		configs/.. \
		internal/..
	@echo "[DONE] Formatting imports..."

gen-mocks:
	@echo "  >  Rebuild Mocking..."

	mockgen -source=pkg/db/mysql.go -destination=gen/mocks/db/mock_db.go -package=db_mocks

	mockgen  --package mockgen -source=internal/loan/services/loan_service.go -destination=gen/mocks/loan/loan_service_mock.go -package=loan_mock
	mockgen  --package mockgen -source=internal/loan/repositories/loan_repository.go -destination=gen/mocks/loan/loan_repository_mock.go -package=loan_mock
	mockgen  --package mockgen -source=internal/loan/repositories/loan_bill_repository.go -destination=gen/mocks/loan/loan_bill_repository_mock.go -package=loan_mock