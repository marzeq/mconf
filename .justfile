default:
  @just --list

builddir := "build"

build:
  mkdir -p {{builddir}}/current-arch
  go build -o {{builddir}}/current-arch/mconf .

build-target OS ARCH:
  mkdir -p {{builddir}}/{{OS}}-{{ARCH}}
  GOOS={{OS}} GOARCH={{ARCH}} go build -o {{builddir}}/{{OS}}-{{ARCH}}/mconf .

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

