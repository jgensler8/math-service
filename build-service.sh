#!/bin/sh

SCRIPT_ROOT=$(dirname $0)

BUILD_DIR=${SCRIPT_ROOT}/docker

docker_prefix="jgensl2/"
docker_version="v0.1.1"
SERVICES="gateway tokenizer addition-operator subtraction-operator multiplication-operator division-operator"

for service in ${SERVICES} ; do
  echo "${service}"
  GOOS=linux go build -o ${BUILD_DIR}/a.out ${SCRIPT_ROOT}/${service}
  docker build -t "${docker_prefix}${service}:${docker_version}" ${BUILD_DIR}
done

for service in ${SERVICES} ; do
  docker push "${docker_prefix}${service}:${docker_version}"
done
