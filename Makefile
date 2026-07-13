GO ?= go
GOVULNCHECK ?= govulncheck
CMD ?=

.PHONY: test verify-dependency-security run vulncheck

test:
	$(GO) test ./...

verify-dependency-security:
	bash scripts/verify-dependency-security.sh



startdev:
	air

startdb:
	./cockroach-v24.1.4.linux-amd64/cockroach start-single-node --insecure=true

sql:
	./cockroach-v24.1.4.linux-amd64/cockroach sql --insecure=true

setup:
	wget https://binaries.cockroachdb.com/cockroach-v24.1.4.linux-amd64.tgz
	tar xvfz cockroach-v24.1.4.linux-amd64.tgz
	$(GO) install github.com/pressly/goose/v3/cmd/goose@latest
	$(GO) install github.com/air-verse/air@latest

vulncheck:
	$(GOVULNCHECK) ./...

run:
	@test -n "$(CMD)" || (echo "usage: make run CMD='go test ./...'" >&2; exit 2)
	$(CMD)
