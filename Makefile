serve-http:
	@go run main.go http
run-cron:
	@go run main.go background
run-worker:
	@go run main.go worker

test:
	./coverage.sh;


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

	# db
	mockgen -source=pkg/db/mysql.go -destination=gen/mocks/db/mock_db.go -package=db_mocks

	# loan
	mockgen  --package mockgen -source=internal/loan/services/loan_service.go -destination=gen/mocks/loan/loan_service_mock.go -package=loan_mock
	mockgen  --package mockgen -source=internal/loan/repositories/loan_repository.go -destination=gen/mocks/loan/loan_repository_mock.go -package=loan_mock
	mockgen  --package mockgen -source=internal/loan/repositories/loan_bill_repository.go -destination=gen/mocks/loan/loan_bill_repository_mock.go -package=loan_mock

	# user
	mockgen  --package mockgen -source=internal/user/services/user_service.go -destination=gen/mocks/user/user_service_mock.go -package=user_mock
	mockgen  --package mockgen -source=internal/user/repositories/user_repository.go -destination=gen/mocks/user/user_repository_mock.go -package=user_mock