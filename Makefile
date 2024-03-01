OUT_FILE_NAME = "api"


run-user:
	cd user-server && go run .

run-goods:
	cd goods-server && go run .

run-stock:
	cd stock-server && go run .

run-user-web:
	cd api-http && make watch-user

watch-user-web:
	cd api-http && make run-user
