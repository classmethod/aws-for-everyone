#!/bin/bash
 
CHANGESET_OPTION="--no-execute-changeset"

if [ $# = 1 ] && [ $1 = "deploy" ]; then
  echo "deploy mode"
  CHANGESET_OPTION=""
fi

# 指定パラメータ
SYSTEM_NAME=handson
TEMPLATE=ecr

# テンプレート実行用パラメータ
CFN_STACK_NAME=${SYSTEM_NAME}-${TEMPLATE}
CFN_TEMPLATE=template/${TEMPLATE}.yml

# テンプレートの実行
aws cloudformation deploy --stack-name ${CFN_STACK_NAME} --template-file ${CFN_TEMPLATE} ${CHANGESET_OPTION} \
  --parameter-overrides \
  SystemName=${SYSTEM_NAME}


# deployモードのとき、ECRの設定変更を実施
if [ $# = 1 ] && [ $1 = "deploy" ]; then
  # 脆弱性自動スキャン設定
  aws ecr put-image-scanning-configuration --repository-name ${SYSTEM_NAME}-ecr-php-app --image-scanning-configuration scanOnPush=true

fi


exit
