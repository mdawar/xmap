.PHONY: test
test:
	go test -cover -race -count 1

.PHONY: benchmark
benchmark:
	go test -bench .
