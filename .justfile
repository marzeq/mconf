default:
  @just --list

builddir := "build"

build:
  mkdir -p {{builddir}}
  go build -o {{builddir}}/mconf .

build-target OS ARCH:
  mkdir -p {{builddir}}
  GOOS={{OS}} GOARCH={{ARCH}} go build -o {{builddir}}/mconf-{{OS}}-{{ARCH}}{{ if OS == "windows" { ".exe" } else { "" } }}

build-all: \
  (build-target "windows" "amd64") \
  (build-target "windows" "arm64") \
  (build-target "linux" "amd64")   \
  (build-target "linux" "arm64")   \
  (build-target "darwin" "amd64")  \
  (build-target "darwin" "arm64")

run *ARGS:
  go run . {{ARGS}}

clean:
  rm -rf {{builddir}}

