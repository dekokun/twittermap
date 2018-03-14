setup-node:
	yarn install
node_modules/%: package.json
	@$(MAKE) setup-node
	@touch $@

.PHONY: build
build: node_modules/.bin/webpack-cli
	node_modules/.bin/webpack-cli

.PHONY: deploy
deploy:
	$(MAKE) build
	aws s3 cp ./public/ s3://twittermap.dekokun.info/2018-03-17/ --acl public-read --recursive --exclude '.DS_Store' --exclude 'node_modules/*' --exclude 'src/*' --exclude gulpfile.js --exclude '.gitignore' --exclude 'Makefile' --exclude 'package*' --exclude 'webpack.config.js' --region ap-northeast-1
.PHONY: dangerous-reset
dangerous-reset:
	aws s3 cp ./tweets.json.default s3://twittermap.dekokun.info/2018-03-17/tweets.json --acl public-read --content-type application/json --region ap-northeast-1 
