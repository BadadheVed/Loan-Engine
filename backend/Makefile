.PHONY: build package deploy clean logs test local


STACK_NAME ?= clickpe-backend-dev
ENVIRONMENT ?= dev
DATABASE_URL ?=
PACKAGE_FILE ?= lambda-functions.zip

build:
	@echo "Building Lambda functions..."
	@chmod +x build-lambda.sh
	@./build-lambda.sh
	@echo "Build complete!"


package:
	@echo "Packaging lambda-functions..."
	@chmod +x package-lambda.sh
	@./package-lambda.sh $(PACKAGE_FILE)
	@echo "Package complete: $(PACKAGE_FILE)"


deploy-guided:
	@echo "Note: Make sure to run 'make build' first"
	cd lambda-functions && sam deploy --guided

# Deploy using SAM
deploy:
	@echo "Deploying with SAM..."
	cd lambda-functions && sam build && sam deploy

# Deploy with DATABASE_URL parameter
deploy-with-db:
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "Error: DATABASE_URL is required. Usage: make deploy-with-db DATABASE_URL='postgres://...'"; \
		exit 1; \
	fi
	cd lambda-functions && sam build && sam deploy --parameter-overrides DatabaseUrl="$(DATABASE_URL)" Environment=$(ENVIRONMENT)

# Start local API
local:
	cd lambda-functions && sam local start-api --env-vars ../env.json

# Test health endpoint locally
test-health:
	curl http://127.0.0.1:3000/api/health


test-upload:
	curl -X POST http://127.0.0.1:3000/api/uploadcsv -F "file=@users.csv"


logs-health:
	sam logs -n HealthFunction --stack-name $(STACK_NAME) --tail

# View logs for uploadcsv function
logs-upload:
	sam logs -n UploadCSVFunction --stack-name $(STACK_NAME) --tail

# Validate SAM template
validate:
	cd lambda-functions && sam validate

# Clean build artifacts
clean:
	rm -rf lambda-functions/.aws-sam
	rm -f lambda-functions/health/bootstrap
	rm -f lambda-functions/uploadcsv/bootstrap
	rm -f $(PACKAGE_FILE)
	rm -f lambda-functions.zip
	@echo "Clean complete!"

# Delete the stack
delete:
	cd lambda-functions && sam delete --stack-name $(STACK_NAME)

# Get API endpoints
endpoints:
	@aws cloudformation describe-stacks \
		--stack-name $(STACK_NAME) \
		--query 'Stacks[0].Outputs' \
		--output table

# Install dependencies
deps:
	cd lambda-functions && go mod download

# Format Go code
fmt:
	cd lambda-functions && go fmt ./...

# Complete workflow: clean, build, package
all: clean build package

# Help
help:
	@echo "Available targets:"
	@echo "  build              - Build Lambda functions"
	@echo "  package            - Create lambda-functions.zip"
	@echo "  deploy-guided      - First time deployment (guided)"
	@echo "  deploy             - Deploy using SAM"
	@echo "  deploy-with-db     - Deploy with DATABASE_URL parameter"
	@echo "  local              - Start local API"
	@echo "  test-health        - Test health endpoint locally"
	@echo "  test-upload        - Test upload endpoint locally"
	@echo "  logs-health        - Tail health function logs"
	@echo "  logs-upload        - Tail upload function logs"
	@echo "  validate           - Validate SAM template"
	@echo "  clean              - Clean build artifacts"
	@echo "  delete             - Delete CloudFormation stack"
	@echo "  endpoints          - Show deployed API endpoints"
	@echo "  deps               - Download Go dependencies"
	@echo "  fmt                - Format Go code"
	@echo "  all                - Clean, build, and package"
	@echo ""
	@echo "Usage examples:"
	@echo "  make build                                      # Build functions"
	@echo "  make package                                    # Create zip"
	@echo "  make deploy-with-db DATABASE_URL='postgres://...' # Deploy"
	@echo "  make all                                        # Complete workflow"
