BUILD_DIR="build"
APP_NAME="dm-db-exporter"
VERSION=$(shell date +'%Y%m%d%H%M%S')
build:
	@go build -o ${BUILD_DIR}/${APP_NAME}__${VERSION}
docker-build:
	@docker build . -t ${APP_NAME}:${VERSION}
