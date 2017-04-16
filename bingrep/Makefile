pkg_base := github.com/s-shin/bintools
pkg := ${pkg_base}/bingrep

build:
	go build ${pkg}

build_debug:
	go build -gcflags "-N -l" ${pkg}

build_test:
	go test -gcflags "-N -l" -c ${pkg}

test:
	go test -v -cover ${pkg}
