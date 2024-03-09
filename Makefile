OUT_FILE_NAME = "api"


run-user:
	cd user-server && go run .

run-goods:
	cd goods-server && go run .

run-stock:
	cd stock-server && go run . -port 10003

run-user-web:
	cd api-http && make watch-user

watch-user-web:
	cd api-http && make run-user

watch-stock:
	@watchexec --restart --ignore docs --exts go make run-stock
