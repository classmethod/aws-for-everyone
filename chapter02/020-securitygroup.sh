#!/bin/bash
 
CHANGESET_OPTION="--no-execute-changeset"

if [ $# = 1 ] && [ $1 = "deploy" ]; then
  echo "deploy mode"
  CHANGESET_OPTION=""
fi

# 指定パラメータ
SYSTEM_NAME=handson
TENPLATE=securitygroup

# 指定必須パラメータ
# 実施環境のグローバルIPアドレスを'curl inet-ip.info'などを利用して、以下に指定してください。
# 例）ADMIN_IPADDRESS=123.123.123.123/32
ADMIN_IPADDRESS=

if [ -z "${ADMIN_IPADDRESS}" ]; then
  echo "シェルの中でADMIN_IPADDRESSを指定してください。"
  exit
fi

# テンプレート実行用パラメータ
CFN_STACK_NAME=${SYSTEM_NAME}-${TENPLATE}
CFN_TEMPLATE=template/${TENPLATE}.yml

# テンプレートの実行
aws cloudformation deploy --stack-name ${CFN_STACK_NAME} --template-file ${CFN_TEMPLATE} ${CHANGESET_OPTION} \
  --parameter-overrides \
  SystemName=${SYSTEM_NAME} \
  AdminIpaddress=${ADMIN_IPADDRESS}

exit
