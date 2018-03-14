CONFIG = config.mk
include $(CONFIG)

OUTPUT_TEMPLATE := output.yml
INPUT_TEMPATE := template.yml

SAMLOCAL := .bin/aws-sam-local


$(OUTPUT_TEMPLATE): $(INPUT_TEMPATE) $(SAMLOCAL) lambda/tweetget.go lambda/s3upload.go
	@$(MAKE) -C lambda build
	$(SAMLOCAL) package \
		--template-file $< \
		--s3-bucket $(CONFIG_CLOUDFORMATION_PACKAGE_S3_BUCKET_NAME) \
		--s3-prefix $(CONFIG_CLOUDFORMATION_PACKAGE_S3_PREFIX) \
		--output-template-file $@

deploy: $(OUTPUT_TEMPLATE) $(SAMLOCAL)
	$(SAMLOCAL) deploy \
		--template-file $< \
		--stack-name $(CONFIG_CLOUDFORMATION_STACK_NAME) \
		--capabilities CAPABILITY_IAM

.bin/%: Makefile
	@$(MAKE) setup-go
	@touch $@

.PHONY: setup-go
setup-go:
	GOBIN=$(abspath .bin) go get -v \
		github.com/awslabs/aws-sam-local

.PHONY: test
test-tweetget: $(SAMLOCAL)
	make -C lambda build
	$(SAMLOCAL) local invoke TweetGetLambda -e event_file.json

test-s3upload: $(SAMLOCAL)
	make -C lambda build
	$(SAMLOCAL) local invoke S3UploadLambda -e event_tweet.json

.PHONY: build
build:
	make -C lambda build
