CONFIG = config.mk
include $(CONFIG)

OUTPUT_TEMPLATE := .output.yml
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
	@$(MAKE) -C web/2018-03-17/ build
	$(SAMLOCAL) deploy \
		--template-file $< \
		--stack-name $(CONFIG_CLOUDFORMATION_STACK_NAME) \
		--parameter-overrides DomainName=$(CONFIG_CLOUDFORMATION_DOMAIN_NAME) TwitterScreenName=$(CONFIG_CLOUDFORMATION_TWITTER_SCREEN_NAME) AcmCertificateArn=$(CONFIG_CLOUDFORMATION_ACM_CERTIFICATE_ARN) \
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
