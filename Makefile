# Пути к прото-файлам
PROTO_PATHS = ./internal/api/grpc
DOCKER_BUF = docker run --rm -w "/work" -v "$$(pwd):/work" anemaloy/go-buf:1.31.0

# Минимальное покрытие кода unit-тестами
UNIT_COVERAGE_MIN = 80
CUR_COVERAGE = 0

.PHONY: all gen_proto unit unit_race lint fmt vet unit_coverage test deps build using

all: using

gen_proto:
	@for CURPATH in $(PROTO_PATHS); do \
		echo "Start processing $$CURPATH..."; \
		rm -Rf $$CURPATH/gen/*; \
		$(MAKE) process_lint_proto_files PROTO_DIR=$$CURPATH/proto; \
		$(MAKE) process_proto_files PROTO_DIR=$$CURPATH/proto TEMPLATE=$$CURPATH/proto/buf.gen.yaml; \
		if [ -d "$$CURPATH/gen" ]; then \
			chown -R $$(id -u):$$(id -g) "$$CURPATH/gen"; \
		fi; \
		echo "Finish processing $$CURPATH..."; \
	done

process_proto_files:
	@if [ ! -d "$$PROTO_DIR" ]; then \
		exit 0; \
	fi; \
	$(DOCKER_BUF) generate $$PROTO_DIR --template $$TEMPLATE

process_lint_proto_files:
	@if [ ! -d "$$PROTO_DIR" ]; then \
		exit 0; \
	fi; \
	$(DOCKER_BUF) lint $$PROTO_DIR

unit:
	@echo "Running unit tests"
	go test ./...

unit_race:
	@echo "Running unit tests with race detection"
	go test -race ./...

lint:
	@echo "Running linter"
	go mod vendor
	$(LINTER_CMD) # Используйте переменную LINTER_CMD для запуска вашего линтера
	rm -Rf vendor

fmt:
	@echo "Running go fmt"
	go fmt ./...

vet:
	@echo "Running go vet"
	go vet ./...

unit_coverage:
	@echo "Running test coverage";
	go test -coverpkg=./... -coverprofile=cover_profile.out.tmp $$(go list ./internal/...)
	< cover_profile.out.tmp grep -v -e "mock" -e "\.pb\.go" -e "\.pb\.validate\.go" > cover_profile.out
	rm cover_profile.out.tmp
	$CUR_COVERAGE=$$(go tool cover -func=cover_profile.out | tail -n 1 | awk '{ print $$3 }')
	rm cover_profile.out
	if [ $(CUR_COVERAGE) -lt $(UNIT_COVERAGE_MIN) ]; then \
		echo "Coverage is not enough $(CUR_COVERAGE) < $(UNIT_COVERAGE_MIN)%"; \
		exit 1; \
	else \
		echo "Coverage is enough $(CUR_COVERAGE) >= $(UNIT_COVERAGE_MIN)%"; \
	fi

test: fmt vet unit unit_race unit_coverage lint

deps:
	@echo "Getting dependencies"
	go get ./...

build: deps
	@echo "Building executable"
	go build ./cmd/template

help:
	@echo "Укажите команду при запуске: make [command]"
	@echo "Список команд:"
	@echo "  gen_proto - генерация прото файлов (для клиентов и сервера)"
	@echo "  unit - запустить unit-тесты"
	@echo "  unit_race - запуск unit тестов с проверкой на data-race"
	@echo "  unit_coverage - запуск unit тестов и проверка покрытия кода тестами"
	@echo "  lint - запустить все линтеры"
	@echo "  fmt - форматирование кода при помощи 'go fmt'"
	@echo "  vet - проверка правильности форматирования кода"
	@echo "  test - запустить все тесты"
	@echo "  deps - подтянуть зависимости"
	@echo "  build - собрать приложение"
