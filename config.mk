# Anything is ok if the bucket name is not exists in the world.
CONFIG_CLOUDFORMATION_PACKAGE_S3_BUCKET_NAME := dekokun-cloudformation-packages
# Anything is ok if the bucket name is not exists in the above bucket.
CONFIG_CLOUDFORMATION_PACKAGE_S3_PREFIX := twittermap
# Anything is ok if the stack name is not exists in your CloudFormation stacks.
CONFIG_CLOUDFORMATION_STACK_NAME := twittermap

CONFIG_CLOUDFORMATION_DOMAIN_NAME := twittermap.dekokun.info
CONFIG_CLOUDFORMATION_TWITTER_SCREEN_NAME := dekokun

# es-east-1 only because cloudfront
CONFIG_CLOUDFORMATION_ACM_CERTIFICATE_ARN := arn:aws:acm:us-east-1:185743233732:certificate/bb90e579-a86b-4056-936a-8bda23020091
