build: build-client build-server

build-client:
	mkdir -p build
	cd client && npm install && npm run-script build
	cp -R client/dist/* ../build/client

build-server:
	mkdir -p build
	go build -tags release -o build/praelatus cmd/praelatus/main.go

clean:
	rm -rf build
	rm -rf client/dist
	rm -rf dist

snapshot: build-client build-server
	if [[ $(git branch | grep "*") != "* master" ]]; then
		echo "ERROR: Must be on master to create a snapshot release"
		exit 1
	fi
	goreleaser --snapshot

package:
	if [[ $(git branch | grep "*") != "* master" ]]; then
		echo "ERROR: Must be on master to create a snapshot release"
		exit 1
	fi
	goreleaser

test-server:
	go test ./...

test-client:
	cd client && npm run-script test
	cd client && npm run-script lint

test: test-server test-client

seeddb:
	go run cmd/praelatus/main.go db seed

cleandb:
	go run cmd/praelatus/main.go db clean --yes

dev-server:
	go run cmd/praelatus/main.go serve --dev-mode

dev-client:
	cd client && npm run dev
