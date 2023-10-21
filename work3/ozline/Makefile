MODULE = github.com/west2-online-reserve/collection-golang/work3

SERVICE_NAME = todolist

.PHONY: target
target:
	sh build.sh
	sh ./output/bootstrap.sh

.PHONY: new
new:
	hz new \
	-module $(MODULE) \
	-service "$(SERVICE_NAME)"
	hz update -idl ./idl/user.thrift
	hz update -idl ./idl/task.thrift

.PHONY: gen
gen:
	hz update -idl ./idl/user.thrift
	hz update -idl ./idl/task.thrift

.PHONY: env
env:
	docker-compose up -d

.PHONY: docker
docker:
	docker build -t west2online .