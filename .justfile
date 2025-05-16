default:
  @just --list

builddir := "build"
projname := "mconf"

build:
  mkdir -p {{builddir}}
  go build -o {{builddir}}/{{projname}} cmd/{{projname}}.go

build-target OS ARCH:
  mkdir -p {{builddir}}
  GOOS={{OS}} GOARCH={{ARCH}} go build -o {{builddir}}/{{projname}}-{{OS}}-{{ARCH}}{{ if OS == "windows" { ".exe" } else { "" } }} cmd/{{projname}}.go

build-all: \
  (build-target "windows" "amd64") \
  (build-target "windows" "arm64") \
  (build-target "linux" "amd64")   \
  (build-target "linux" "arm64")   \
  (build-target "darwin" "amd64")  \
  (build-target "darwin" "arm64")

run *ARGS:
  go run cmd/{{projname}}.go {{ARGS}}

clean:
  rm -rf {{builddir}}

