#!/bin/bash
 
CHANGESET_OPTION="--no-execute-changeset"

if [ $# = 1 ] && [ $1 = "deploy" ]; then
  echo "deploy mode"
  CHANGESET_OPTION=""
fi

# 指定パラメータ
SYSTEM_NAME=handson
TENPLATE=app-alb-fargate

# テンプレート実行用パラメータ
CFN_STACK_NAME=${SYSTEM_NAME}-${TENPLATE}
CFN_TEMPLATE=template/${TENPLATE}.yml
ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text)
IMAGE_URI=${ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/${SYSTEM_NAME}-ecr-php-app

# テンプレートの実行
aws cloudformation deploy --stack-name ${CFN_STACK_NAME} --template-file ${CFN_TEMPLATE} ${CHANGESET_OPTION} \
  --capabilities CAPABILITY_IAM \
  --parameter-overrides \
  SystemName=${SYSTEM_NAME} \
  ImageUri=${IMAGE_URI}

exit
