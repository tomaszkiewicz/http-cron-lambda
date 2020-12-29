.DEFAULT_GOAL := publish

build:
	@echo "==> Building"
	mkdir -p dist
	GOOS=linux go build -o dist/main ./cmd/lambda

package: build
	@echo "==> Packaging"
	rm -f dist/lambda.zip
	zip -j -r dist/lambda.zip dist/main

publish: package
	@echo "==> Publishing"
	#AWS_PROFILE=management
	aws lambda update-function-code --function-name http-cron --zip-file fileb://dist/lambda.zip
