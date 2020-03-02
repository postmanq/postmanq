#!/usr/bin/env bash

PROJECT_DIR=$(pwd)
MODULE_DIR=${PROJECT_DIR}/module
OUT_DIR=${PROJECT_DIR}/out
OUT_MODULE_DIR=${OUT_DIR}/module

rm -rf ${OUT_DIR}
mkdir -p ${OUT_MODULE_DIR}

for DIR in $(ls -d ${MODULE_DIR}/*/); do
  MODULE_NAME=$(basename ${DIR})
  go build -buildmode=plugin -o ${OUT_MODULE_DIR}/${MODULE_NAME}.so ${MODULE_DIR}/${MODULE_NAME}/module.go
done;

go build -o ${OUT_DIR}/postmanq ${PROJECT_DIR}/cmd/postmanq/postmanq.go

${OUT_DIR}/postmanq -c ${PROJECT_DIR}/config.example.yml