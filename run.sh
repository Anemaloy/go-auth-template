#!/bin/sh -e
PROTO_PATHS="internal/api/grpc"
UNIT_COVERAGE_MIN=4

# Запуск buf
run_buf(){
  docker run --rm -w "/work" -v "$(pwd):/work" anemaloy/go-buf:1.31.0 $@
}

# Выполняет команду от имени root
run_as_root(){
  docker run --rm -w "/work" -v "$(pwd):/work" anemaloy/go-alpine:3 $@
}

# Обрабатывает прото файлы
process_proto_files(){
  local COMMAND="$1"
  local PROTO_DIR="$2"

  if [ ! -d "$PROTO_DIR" ]; then
    return 0
  fi

  run_buf $@
}

# Генерация прото файлов
gen_proto(){
  for CURPATH in ${PROTO_PATHS}; do
    echo "start process $CURPATH..."

    rm -Rf $CURPATH/gen/*
    process_proto_files lint "$CURPATH/proto"
    process_proto_files generate "$CURPATH/proto" --template "$CURPATH/proto/buf.gen.yaml"

    if [ -d "$CURPATH/gen" ]; then
          run_as_root chown -R "$(id -u)":"$(id -g)" "/work/$CURPATH/gen"
    fi

    echo "finish process $CURPATH..."
  done
}

# Запуск unit-тестов
unit(){
  echo "run unit tests"
  go test ./...
}

unit_race() {
  echo "run unit tests with race test"
  go test -race ./...
}

# Запуск go-lint
lint(){
  echo "run linter"
  go mod vendor
  docker run --rm -v $(pwd):/work:ro -w /work golangci/golangci-lint:latest golangci-lint run -v
  rm -Rf vendor
}

fmt() {
  echo "run go fmt"
  go fmt ./...
}

vet() {
  echo "run go vet"
  go vet ./...
}

unit_coverage() {
  echo "run test coverage"
  go test -coverpkg=./... -coverprofile=cover_profile.out.tmp $(go list ./internal/...)
  # remove generated code and mocks from coverage
  < cover_profile.out.tmp grep -v -e "mock" -e "\.pb\.go" -e "\.pb\.validate\.go" > cover_profile.out
  rm cover_profile.out.tmp
  CUR_COVERAGE=$( go tool cover -func=cover_profile.out | tail -n 1 | awk '{ print $3 }' | sed -e 's/^\([0-9]*\).*$/\1/g' )
  rm cover_profile.out
  if [ "$CUR_COVERAGE" -lt $UNIT_COVERAGE_MIN ]
  then
    echo "coverage is not enough $CUR_COVERAGE < $UNIT_COVERAGE_MIN"
    return 1
  else
    echo "coverage is enough $CUR_COVERAGE >= $UNIT_COVERAGE_MIN"
  fi
}

# Запуск всех тестов
test(){
  fmt
  vet
  unit
  unit_race
  unit_coverage
  lint
}

# Подтянуть зависимости
deps(){
  go get ./...
}

# Собрать исполняемый файл
build(){
  deps
  go build ./cmd/template
}

# Добавьте сюда список команд
help(){
  echo "Укажите команду при запуске: ./run.sh [command]"
  echo "Список команд:"
  echo "  unit - запустить unit-тесты"
  echo "  unit_race - запуск unit тестов с проверкой на data-race"
  echo "  unit_coverage - запуск unit тестов и проверка покрытия кода тестами"
  echo "  lint - запустить все линтеры"
  echo "  test - запустить все тесты"
  echo "  deps - подтянуть зависимости"
  echo "  build - собрать приложение"
  echo "  fmt - форматирование кода при помощи 'go fmt'"
  echo "  vet - проверка правильности форматирования кода"
  echo "  gen_proto - генерация прото файлов (для клиентов и сервера)"
}

############### НЕ МЕНЯЙТЕ КОД НИЖЕ ЭТОЙ СТРОКИ #################

command="$1"
if [ -z "$command" ]
then
 help
 exit 0;
else
 $command $@
fi
