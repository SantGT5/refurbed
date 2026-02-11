include $(wildcard task/*.mk)

.DEFAULT_GOAL := help

urls: ## Show the urls to the running applications
	@echo "*------"
	@echo "* Assignment Backend"
	@echo "*"
	@echo "* Backend: http://localhost:8080"
	@echo "* Frontend: http://localhost:5173"
	@echo "*"
	@echo "*------\n"
.PHONY: urls


build: ## Build the application
	@go build -o tmp/main main.go
.PHONY: build

start: build urls ## Run the application
	@./tmp/main
.PHONY: start

format-check: ## Check code formatting
	@gofmt -l .
	@if [ -n "$$(gofmt -l .)" ]; then echo "Code is not formatted. Run 'make format' to format.\n"; exit 1; fi
.PHONY: quality

format: ## Format the code
	@gofmt -w .
.PHONY: format
