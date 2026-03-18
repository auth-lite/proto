.PHONY: clean generate update tidy

clean:
	@rm -rf gen/go

generate:
	@buf generate

update:
	@buf dep update

tidy:
	@buf format -w
	@buf lint