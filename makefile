.PHONY: build zip deploy invoke logs clean

FUNCTION_NAME = connect-route-lambda
HANDLER_NAME = bootstrap
ZIP_NAME = connect-route-lambda.zip
SOURCE_FILES = $(shell find . -name '*.go')
LOG_GROUP_NAME = /aws/lambda/$(FUNCTION_NAME)
INPUT_FILE = input.json
OUTPUT_FILE = output.json

all: build zip deploy invoke clean

update: build zip deploy clean

# Build the Go binary
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda -o $(HANDLER_NAME)

# Zip the binary
zip:
	zip $(ZIP_NAME) $(HANDLER_NAME)

# Deploy the zipped binary to Lambda
deploy:
	aws lambda update-function-code --function-name $(FUNCTION_NAME) --no-cli-pager --zip-file fileb://$(ZIP_NAME)

# Invoke the Lambda function
invoke:
	@echo "Invoking Lambda function..."
	aws lambda invoke --function-name $(FUNCTION_NAME) --payload file://$(INPUT_FILE) $(OUTPUT_FILE)
	@echo "Lambda Response saved in output.json"

# Fetch the latest logs
logs:
	@echo "Fetching logs..."
	@latest_stream=$$(aws logs describe-log-streams --log-group-name $(LOG_GROUP_NAME) --order-by LastEventTime --descending --max-items 1 --query 'logStreams[0].logStreamName' --output text) && \
	aws logs get-log-events --log-group-name $(LOG_GROUP_NAME) --log-stream-name $$latest_stream

clean:
	rm -f $(HANDLER_NAME)
	rm -f $(ZIP_NAME)
