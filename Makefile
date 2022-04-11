PROJECT_NAME=storage
GIT_COMMIT=$(shell git rev-list -1 HEAD)

PKGDIR=pkg

TARGET_BIN=storage
TARGET_PACKAGE=*

HOST_PORT=8081
CONTAINER_PORT=8081

TAG=latest


GOOS=$(shell uname | tr '[:upper:]' '[:lower:]')
GOARCH=amd64
CGO=0

ENV_VARS=GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO}

# local install
install: 
	go mod download

${TARGET_BIN}: install
	env ${ENV_VARS} go build cmd/${TARGET_BIN}.go 

run:
	./${TARGET_BIN}

# test:
# 	env ${ENV_VARS} go test -ldflags "	-X github.com/barbabjetolov/endocode-test/http-service/pkg/utilities.ProjectName=${PROJECT_NAME} \
# 										-X github.com/barbabjetolov/endocode-test/http-service/pkg/utilities.GitCommit=${GIT_COMMIT} 	\
# 										-w -s" \
# 										./${PKGDIR}/${TARGET_PACKAGE} 


# all: test ${TARGET_BIN}


#docker
docker-build:
	docker build -t ${TARGET_BIN}:${TAG} .

docker-run:
	docker run -d --name ${TARGET_BIN} -e LISTENING_PORT=${CONTAINER_PORT} -dp ${HOST_PORT}:${CONTAINER_PORT} -it ${TARGET_BIN}:${TAG}

docker-clean:
	-docker image rm ${TARGET_BIN}:${TAG}

docker: docker-build docker-run