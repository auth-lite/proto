.PHONY: proto clean tidy

export GOPRIVATE=github.com/kperreau

clean:
	@rm -rf gen/go

generate:
	@buf generate

update:
	@buf dep update

tidy:
	@buf format -w
	@buf lint