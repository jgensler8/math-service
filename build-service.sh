#!/bin/sh

SCRIPT_ROOT=$(dirname $0)

BUILD_DIR=${SCRIPT_ROOT}/docker

docker_prefix=""
SERVICES="gateway tokenizer addition-operator subtraction-operator multiplication-operator division-operator"

for service in ${SERVICES} ; do
  echo "${service}"
  GOOS=linux go build -o ${BUILD_DIR}/a.out ${SCRIPT_ROOT}/${service}
  docker build -t "${docker_prefix}${service}" ${BUILD_DIR}
done
