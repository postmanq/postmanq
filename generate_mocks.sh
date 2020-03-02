#!/usr/bin/env bash
GOBIN=~/GoglandProjects/bin

${GOBIN}/mockery -dir=module/config/service -all -case=snake -outpkg=service_mock -output=module/config/service/mock