# go test -run TestLoginSuccess  ./src/handlers

test-unit: 
	go test ./src/util -coverprofile=.build/test-coverage/util-cover.out

test-intergration: 
	go test ./src/repo -coverprofile=.build/test-coverage/repo-cover.out
	go test ./src/handlers -coverprofile=.build/test-coverage/handler-cover.out
	go test ./src/middleware -coverprofile=.build/test-coverage/middleware-cover.out

test: test-unit test-intergration

cover-util:
	go tool cover -html=.build/test-coverage/util-cover.out

cover-repo:
	go tool cover -html=.build/test-coverage/repo-cover.out

cover-handler:
	go tool cover -html=.build/test-coverage/handler-cover.out

cover-middleware:
	go tool cover -html=.build/test-coverage/middleware-cover.out

deploy-prod: test
	sh deploy_prod.sh 

deploy-stage: test
	sh deploy_stage.sh 

run:
	go run main.go

# .PHONY : test-intergration test-unit test
