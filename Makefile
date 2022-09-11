version = "v1.5.0"
srcPath = "tmp/$(version)/src"

install-srt:
	rm -rf $(srcPath)
	mkdir -p $(srcPath)
	git clone https://github.com/Haivision/srt $(srcPath)
	cd $(srcPath) && git checkout $(version)
	cd $(srcPath) && ./configure --prefix=.. $(configure)
	cd $(srcPath) && make
	cd $(srcPath) && make install

generate:
	go run internal/cmd/options/main.go
	go run internal/cmd/stats/main.go
	go run internal/cmd/wrap/main.go

test-coverage:
	go test -coverprofile cover.out github.com/asticode/go-astisrt/pkg
	go tool cover -html=cover.out
