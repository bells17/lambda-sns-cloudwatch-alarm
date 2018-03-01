FILE_NAME=lambda-sns-cloudwatch-alarm

dep:
	dep ensure

build:
	GOOS=linux GOARCH=amd64 go build -o ${FILE_NAME}

pack: build
	zip ${FILE_NAME}.zip ${FILE_NAME}

clean:
	rm -f ${FILE_NAME} ${FILE_NAME}.zip

fmt:
	go fmt
	goimports -w ./*.go

publish:
	aws lambda update-function-code --function-name ${FUNCTION_NAME} --zip-file fileb://./${FILE_NAME}.zip --publish $(PUBLISH_OPTIONS)
