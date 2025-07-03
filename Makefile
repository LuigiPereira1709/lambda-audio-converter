# Load environment variables
include .env
export $(shell sed 's/=.*//' .env)

FUNCTION_NAME := ffmpeg-lambda
FUNCTION_ZIP := ./lambda.zip
LAYER_NAME := ffmpeg-layer
LAYER_ZIP := ./ffmpeg.zip
RUNTIME := provided.al2023 
HANDLER := bootstrap 
ARCHITECTURE := arm64
MEMORY_SIZE := 640
TIMEOUT := 480
LAMBDA_ENV_PATH := ./.lambda_env.json

.DEFAULT_GOAL := help

# Ensure logs directory exists
.PHONY: init-logs
init-logs:
	@mkdir -p logs

# Get the latest layer ARN
.PHONY: get-layer-arn
get-layer-arn:
	@aws lambda list-layer-versions \
		--layer-name $(LAYER_NAME) \
		--query "LayerVersions[0].LayerVersionArn" \
		--output text

# Build the lambda function build-lambda:
.PHONY: build-lambda
build-lambda:
	@./scripts/lambda_build.sh

# Build the FFmpeg layer
.PHONY: build-ffmpeg-layer
build-ffmpeg-layer:
	@./scripts/ffmpeg_layer_build.sh

# Deploy the FFmpeg layer
.PHONY: deploy-ffmpeg-layer
deploy-ffmpeg-layer: init-logs
	@echo "Deploying the FFmpeg layer..."
	@if [ ! -f "$(LAYER_ZIP)" ]; then \
		echo "Error: Layer zip file '$(LAYER_ZIP)' does not exist."; \
		exit 1; \
	fi
	@aws lambda publish-layer-version \
		--layer-name "$(LAYER_NAME)" \
		--description "FFmpeg layer for Lambda" \
		--zip-file fileb://$(LAYER_ZIP) \
		--compatible-architectures $(ARCHITECTURE) \
		--compatible-runtimes $(RUNTIME) > logs/publish-layer-version.log 2>&1
	@echo "FFmpeg layer deployed successfully!"

# Deploy the lambda function
.PHONY: deploy-lambda
deploy-lambda: init-logs
	@echo "Deploying the lambda function..."
	@if [ ! -f "$(FUNCTION_ZIP)" ]; then \
		echo "Error: Function zip file '$(FUNCTION_ZIP)' does not exist."; \
		exit 1; \
	fi
	@set -e; \
	echo "--> Fetching Layer ARN..."; \
	LAYER_ARN=$$(make --no-print-directory get-layer-arn); \
	echo "--> Layer ARN: $$LAYER_ARN"; \
	if aws lambda get-function --function-name $(FUNCTION_NAME) > /dev/null 2>&1; then \
		echo "--> Function $(FUNCTION_NAME) exists, updating..."; \
		echo "--> Updating function code..."; \
		aws lambda update-function-code \
			--function-name $(FUNCTION_NAME) \
			--zip-file fileb://$(FUNCTION_ZIP) \
			--publish >> logs/update-function-code.log 2>&1; \
		echo "--> Waiting for function to be ready..."; \
		aws lambda wait function-updated \
			--function-name $(FUNCTION_NAME) > logs/wait-updated.log 2>&1; \
		echo "--> Updating function configuration..."; \
		aws lambda update-function-configuration \
			--function-name $(FUNCTION_NAME) \
			--role $(LAMBDA_ROLE) \
			--handler $(HANDLER) \
			--memory-size $(MEMORY_SIZE) \
			--timeout $(TIMEOUT) \
			--environment file://$(LAMBDA_ENV_PATH) \
			--layers $$LAYER_ARN >> logs/update-config.log 2>&1; \
	else \
		echo "--> Function $(FUNCTION_NAME) does not exist, creating..."; \
		aws lambda create-function \
			--function-name $(FUNCTION_NAME) \
			--runtime $(RUNTIME) \
			--zip-file fileb://$(FUNCTION_ZIP) \
			--handler $(HANDLER) \
			--layers $$LAYER_ARN \
			--architectures $(ARCHITECTURE) \
			--role $(LAMBDA_ROLE) \
			--memory-size $(MEMORY_SIZE) \
			--timeout $(TIMEOUT) \
			--environment file://$(LAMBDA_ENV_PATH) \
			--description "Lambda function for FFmpeg processing" >> logs/create-function.log 2>&1; \
	fi
	@echo "Lambda function deployed successfully!"

# Sloth mode: build (lambda, FFmpeg layer) and deploy
.PHONY: sloth
sloth: build-lambda build-ffmpeg-layer deploy-ffmpeg-layer deploy-lambda
	@echo "🦥 Yeah, I am a sloth!"

# Help command
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build-lambda           Build the Lambda Go binary"
	@echo "  build-ffmpeg-layer     Build the FFmpeg Lambda layer"
	@echo "  deploy-ffmpeg-layer    Deploy the FFmpeg layer to AWS Lambda"
	@echo "  deploy-lambda          Deploy the Lambda function"
	@echo "  get-layer-arn          Print the latest Layer ARN"
	@echo "  sloth                  Build and deploy all (lazy mode 🦥)"
