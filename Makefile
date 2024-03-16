rinha:
	@echo building container image
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/rinha-app ./cmd/*.go && docker build -t moaabb/rinha-go-2024-q1:1.0 . 