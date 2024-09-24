version = "v1.5.3"
srcPath = "tmp/$(version)/src"
currentDir := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

install-srt:
	rm -rf $(srcPath)
	mkdir -p $(srcPath)
	# cd $(srcPath) is necessary for windows build since otherwise git doesn't clone in the proper dir
	cd $(srcPath) && git clone https://github.com/Haivision/srt .
	cd $(srcPath) && git checkout $(version)
	cd $(srcPath) && ./configure --prefix=.. $(configure)
	cd $(srcPath) && make
	cd $(srcPath) && make install

generate:
	go run internal/cmd/generate/options/main.go
	go run internal/cmd/generate/static_consts/main.go
	go run internal/cmd/generate/stats/main.go
	go run internal/cmd/generate/wrap/main.go

test-coverage:
	go test -coverprofile cover.out github.com/asticode/go-astisrt/pkg
	go tool cover -html=cover.out

test-linux-build:
	cd testdata/linux && docker build -t astisrt-test-linux .

test-linux:
	docker run -v ${currentDir}/testdata/linux/gocache:/opt/astisrt/tmp/linux/gocache -v ${currentDir}/testdata/linux/gopath:/opt/astisrt/tmp/linux/gopath -v ${currentDir}:/opt/astisrt/tmp/linux/gopath/src/github.com/asticode/go-astisrt astisrt-test-linux