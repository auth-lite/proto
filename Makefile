.PHONY: clean gen update tidy

clean:
	@rm -rf gen/go

gen:
	@buf gen

update:
	@buf dep update

tidy:
	@buf format -w
	@buf lint