OUT_FILE_NAME = "api"

run-user:
	@go run ./user-web/main.go

watch-user:
	@watchexec --restart --ignore docs --exts go make run-user

kill_and_run:
	@killport 11001

release:
	@go build -o ${OUT_FILE_NAME}

start:
	@yarn --cwd web-ui dev

build-web:
	@yarn --cwd web-ui build


