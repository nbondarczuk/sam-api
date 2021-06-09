TARGET = sam-api
VERSION = $(shell git describe --tags)
BUILD = $(shell date +"%F_%T_%Z")
LEVEL = $(shell git log --pretty=format:"%H" --name-status HEAD^..HEAD | head -1)
IMAGE = registry1.corpo.t-mobile.pl/sam/sam-api:$(VERSION)
IMAGE_BASE = registry1.corpo.t-mobile.pl/sam/centos-oracle-base
DOCKERFILE = Dockerfile.centos
DOCKERFILE_BASE = Dockerfile-base.centos

DEPS = "github.com/codegangsta/negroni" \
"github.com/dgrijalva/jwt-go" \
"github.com/dgrijalva/jwt-go/request" \
"github.com/gorilla/mux" \
"gopkg.in/goracle.v2" \
"golang.org/x/crypto/bcrypt" \
"github.com/unrolled/render" \
"github.com/Sirupsen/logrus" \
"gopkg.in/gorp.v2" \
"github.com/go-chi/chi" \
"gopkg.in/ldap.v2" \
"github.com/jtblin/go-ldap-client"

LDFLAGS = "-X main.version=$(VERSION) -X main.build=$(BUILD) -X main.level=$(LEVEL)"
STATIC_BUILD_PREFIX = "CGO_ENABLED=1 GOOS=linux GOARCH=amd64"
STATIC_LDFLAGS = "-w -extldflags '-static -I$(GOPATH)/src/gopkg.in/goracle.v2/odpi/include -I$(GOPATH)/src/gopkg.in/goracle.v2/odpi/src -I$(GOPATH)/src/gopkg.in/goracle.v2/odpi/embed -ldl' "

all: build

build:
	go build -ldflags=$(LDFLAGS) -o $(TARGET) $(TARGET).go

build-static:
	$(STATICS_BUILD_PREFIX) go build -tags netgo -ldflags=$(LDFLAGS) -ldflags=$(STATIC_LDFLAGS) -o $(TARGET) $(TARGET).go

run:
	go run sam-api.go

run-local:
	cp config/config.json.XE config.json
	go run sam-api.go

install:
	go install

deps:
	go get $(DEPS)
	sudo apt install jsonlint curl openssl

keys:
	openssl genrsa -out keys/app.rsa
	openssl rsa -in keys/app.rsa -pubout > keys/app.rsa.pub

sql: sql-xe

sql-xe:
	(cd sql; bash install-xe.sh)

sql-T17:
	(cd sql; bash install-T17.sh)

docker: build keys
	echo $(VERSION) > VERSION
	cp config/config.json.T17 config.json
	docker build -f $(DOCKERFILE) -t $(IMAGE) .

docker-base:
	docker build -f $(DOCKERFILE_BASE) -t $(IMAGE_BASE) .

docker-run:
	docker run -it --rm --network host --name $(TARGET) $(IMAGE)

docker-run-detached:
	docker run -d --network host --name $(TARGET) $(IMAGE)

docker-push:
	docker push $(IMAGE)

docker-base-push:
	docker push $(IMAGE_BASE)

swagger-ui:
	docker run --rm --network host -e SWAGGER_JSON=/foo/swagger.yaml -v $(PWD)/swagger:/foo swaggerapi/swagger-ui

swagger-editor:
	docker run --rm --network host -e SWAGGER_JSON=/foo/swagger.yaml -v $(PWD)/swagger:/foo swaggerapi/swagger-editor

clean:
	go clean
	rm -f *~ */*~ $(TARGET) keys/* test/results/*.* test/runtest/results/*.* test/runtest/logs/*.*

tar: clean
	rm -f  ~/sam-api.tar.gz
	tar -cvf ~/sam-api.tar *
	gzip  -f ~/sam-api.tar

test: unittest-xe

unittest-t17: keys
	cp config/config.json.T17 config.json
	CONFIG=${PWD}/config.json RUNPATH=${PWD} KEYPATH=${PWD}/keys go test -v -cover ./...

start-xe:
	docker run -d -v /var/oracle-xe-data:/u01/app/oracle -p 1521:1521 -e ORACLE_ALLOW_REMOTE=true -e ORACLE_PASSWORD=oracle -e RELAX_SECURITY=1 epiclabs/docker-oracle-xe-11g

unittest-xe: keys
	cp config/config.json.XE config.json
	CONFIG=${PWD}/config.json RUNPATH=${PWD} KEYPATH=${PWD}/keys go test -v -cover ./...

unittest-xe-1: keys
	cp config/config.json.XE config.json
	CONFIG=${PWD}/config.json RUNPATH=${PWD} KEYPATH=${PWD}/keys go test -run TestAccountCreate ./controllers/unittest/ -v -count 1

unittest-xe-2: keys
	cp config/config.json.XE config.json
	CONFIG=${PWD}/config.json RUNPATH=${PWD} KEYPATH=${PWD}/keys go test -run TestDictionaryAccountSapCreateExcel ./controllers/unittest/ -v -count 1

unittest-xe-3: keys
	cp config/config.json.XE config.json
	CONFIG=${PWD}/config.json RUNPATH=${PWD} KEYPATH=${PWD}/keys go test -run TestAccountReadAllAsCsv ./controllers/unittest/ -v -count 1

.PHONY: test clean docker keys run build all deps test unittest runtest test-login test-segment-create test-segment-read test-segment-delete install sql
