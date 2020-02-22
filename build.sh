#!/usr/bin/env bash

PROJECT_DIR=$(pwd)
PLUGIN_DIR=${PROJECT_DIR}/plugin
OUT_DIR=${PROJECT_DIR}/out
OUT_PLUGIN_DIR=${OUT_DIR}/plugin

rm -rf ${OUT_PLUGIN_DIR}
mkdir -p ${OUT_PLUGIN_DIR}

for DIR in $(ls -d ${PLUGIN_DIR}/*/); do
  PLUGIN_NAME=$(basename ${DIR})
  go build -buildmode=plugin -o ${OUT_PLUGIN_DIR}/${PLUGIN_NAME}.so ${PLUGIN_DIR}/${PLUGIN_NAME}/plugin.go
done;