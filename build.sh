#!/bin/bash
set -e

# Check go install
if [ "$(which go)" == "" ]; then
	echo "error: Go is not installed. Please download and follow installation instructions at https://golang.org/dl to continue."
	exit 1
fi


cd $(dirname "${BASH_SOURCE[0]}")
OD="$(pwd)"


if [ "$1" == "vendor" ]; then
	pkg="$2"
	if [ "$pkg" == "" ]; then
		echo "no package specified"
		exit
	fi
	if [ ! -d "$GOPATH/src/$pkg" ]; then
		echo "invalid package"
		exit
	fi
	rm -rf vendor/$pkg/
	mkdir -p vendor/$pkg/
	cp -rf $GOPATH/src/$pkg/* vendor/$pkg/
	rm -rf vendor/$pkg/.git
	exit
fi

# temp directory for storing isolated environment.
TMP="$(mktemp -d -t juno-server.XXXX)"
function rmtemp {
	rm -rf "$TMP"
}
trap rmtemp EXIT

if [ "$NOCOPY" != "1" ]; then
	# copy all files to an isloated directory.
	WD="$TMP/src/github.com/junostorage"
	export GOPATH="$TMP"

	for file in `find . -type f`; do
		# TODO: use .gitignore to ignore, or possibly just use git to determine the file list.
		if [[ "$file" != "." && "$file" != ./.git* && "$file" != ./data* && "$file" != ./juno-server* ]]; then

			mkdir -p "$WD/$(dirname "${file}")"
			cp -P "$file" "$WD/$(dirname "${file}")"
		fi
	done
	cd $WD
fi

# test if requested
if [ "$1" == "test" ]; then
	$OD/juno-server -p 6380  &
	PID=$!
	function testend {
		kill $PID &
	}
	trap testend EXIT
	go test $(go list ./... | grep -v /vendor/)
fi

# cover if requested
if [ "$1" == "cover" ]; then
	$OD/juno-server  -p 6380  &
	PID=$!
	function testend {
		kill $PID &
	}
	trap testend EXIT
	go test -cover $(go list ./... | grep -v /vendor/)
fi

# build and store objects into original directory.
go build -ldflags "$LDFLAGS" -o "$OD/juno-server" cmd/juno-server/*.go
