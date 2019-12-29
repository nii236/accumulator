# Go parameters
include .env
export $(shell sed 's/=.*//' .env)

all: clean prepare deps build-server build-web copy
prepare: 
	mkdir deploy
	cd deploy && mkdir bin config web
build-server: 
	go generate
	go run cmd/admin/main.go -db-drop
	go run cmd/admin/main.go -db-migrate
	go generate
	go build -o deploy/bin/accumulator cmd/accumulator/main.go
	go build -o deploy/bin/admin cmd/admin/main.go
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
	rsync -avz -e 'ssh -p $(PROD_PORT)' ./deploy/ $(PROD_USER)@$(PROD_HOST):$(PROD_PATH)
	ssh -p $(PROD_PORT) $(PROD_USER)@$(PROD_HOST) sudo systemctl daemon-reload
	ssh -p $(PROD_PORT) $(PROD_USER)@$(PROD_HOST) sudo systemctl restart accumulator 
deploy-prod-frontend:
	rsync -avz -e 'ssh -p $(PROD_PORT)' ./deploy/web/ $(PROD_USER)@$(PROD_HOST):$(PROD_PATH)/web