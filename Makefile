OUT_FILE_NAME = "api"


run-user:
	cd user-server && go run . -port 10001

run-goods:
	cd goods-server && go run . -port 10002

run-goods2:
	cd goods_srv && go run . -port 10002

run-stock:
	cd stock-server && go run . -port 10003

run-order:
	cd order-server && go run . -port 10004

run-user-web:
	cd api-http && make watch-user

watch-user-web:
	cd api-http && make run-user

watch-goods2:
	@watchexec --restart --ignore docs --exts go make run-goods2

watch-goods:
	@watchexec --restart --ignore docs --exts go make run-goods

watch-stock:
	@watchexec --restart --ignore docs --exts go make run-stock

watch-order:
	@watchexec --restart --ignore docs --exts go make run-order
