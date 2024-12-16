serve-http:
	@go run main.go http

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