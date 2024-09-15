


startdev:
	air

startdb:
	./cockroach-v24.1.4.linux-amd64/cockroach start-single-node --insecure=true

setup:
	wget https://binaries.cockroachdb.com/cockroach-v24.1.4.linux-amd64.tgz
	tar xvfz cockroach-v24.1.4.linux-amd64.tgz
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/air-verse/air@latest
