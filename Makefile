generate:
	go run internal/cmd/generate/options/main.go
	go run internal/cmd/generate/static_consts/main.go
	go run internal/cmd/generate/stats/main.go
	go run internal/cmd/generate/wrap/main.go