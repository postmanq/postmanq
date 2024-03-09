FROM golang:alpine

ARG BUILD_APP_NAME

ENV PROJECT_DIR=/go
ENV DIST_DIR=${PROJECT_DIR}/dist
ENV GOOS=linux
ENV GOARCH=arm64

RUN apk add --no-cache tzdata gcc musl-dev make binutils-gold

COPY . ${PROJECT_DIR}

RUN make build_postmanq && make build_postmanq build_plugins

ENTRYPOINT /bin/ash -c "chmod +x ${DIST_DIR}/postmanq-${GOOS}-${GOARCH} && ${DIST_DIR}/postmanq-${GOOS}-${GOARCH} -c ${PROJECT_DIR}/configs/config.yml"
