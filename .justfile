default:
  @just --list

build:
  mkdir -p build 
  go build -o build/mconf .

run *ARGS:
  go run . {{ARGS}}

clean:
  rm -rf build
