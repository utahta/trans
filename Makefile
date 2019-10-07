export GO111MODULE=on

mod:
	go mod download

build:
	go build -o ./bin/trans ./cmd/trans

package:
	./scripts/package.sh

release: package
	./scripts/release.sh

