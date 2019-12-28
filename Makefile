# Go parameters
include .env
export $(shell sed 's/=.*//' .env)

all: clean prepare deps build-server build-web copy
prepare: 
	mkdir deploy
	cd deploy && mkdir bin config email migrations init web certs assets
build-server: 
	go generate
	go run cmd/accumulator/main.go -db-migrate
	go generate
	go build -o deploy/bin/accumulator cmd/server/main.go
build-web:
	cd web && npm install
	cd web && npm run build
clean:
	rm -rf deploy
copy:
	cp -r config deploy/
	cp -r web/dist/. deploy/web
deps:
	go mod download
deploy-prod-full:
	rsync -avz ./deploy/ $(PROD_USER)@$(PROD_HOST):$(PROD_PATH)
deploy-prod-frontend:
	rsync -avz ./deploy/web/ $(PROD_USER)@$(PROD_HOST):$(PROD_PATH)/web