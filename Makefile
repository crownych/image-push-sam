#!/usr/bin/env make

#PROFILE=default
PROFILE=dev36-sso
REGION=ap-northeast-1

STACK_NAME=image-push-sam-stack

.PHONY: deploy
deploy: validate build
	sam deploy \
		--profile ${PROFILE} \
		--region ${REGION} \
		--stack-name ${STACK_NAME} \
		--resolve-s3 \
		--resolve-image-repos \
		--capabilities CAPABILITY_IAM \
		--confirm-changeset \
		--debug
	aws cloudformation describe-stacks \
 		--stack-name ${STACK_NAME} \
		--profile ${PROFILE} \
		--region ${REGION} \
		--query 'Stacks[].Outputs'

.PHONY: delete
delete:
	sam delete \
    	--profile ${PROFILE} \
    	--region ${REGION} \
    	--stack-name ${STACK_NAME}

.PHONY: validate
validate:
	sam validate \
		--profile ${PROFILE} \
		--region ${REGION} \
		--debug

.PHONY: build
build:
	sam build \
        --profile ${PROFILE} \
        --region ${REGION} \
        --debug
