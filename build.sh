#!/bin/bash
set -e


# Check go install
if [ "$(which go)" == "" ]; then
	echo "error: Go is not installed. Please download and follow installation instructions at https://golang.org/dl to continue."
	exit 1
fi

# Check go version
GOVERS="$(go version | cut -d " " -f 3)"
if [ "$GOVERS" != "devel" ]; then
	vercomp () {
		if [[ $1 == $2 ]]
		then
			echo "0"
			return
		fi
		local IFS=.
		local i ver1=($1) ver2=($2)
		# fill empty fields in ver1 with zeros
		for ((i=${#ver1[@]}; i<${#ver2[@]}; i++))
		do
			ver1[i]=0
		done
		for ((i=0; i<${#ver1[@]}; i++))
		do
			if [[ -z ${ver2[i]} ]]
			then
				# fill empty fields in ver2 with zeros
				ver2[i]=0
			fi
			if ((10#${ver1[i]} > 10#${ver2[i]}))
			then
				echo "1"
				return
			fi
			if ((10#${ver1[i]} < 10#${ver2[i]}))
			then
				echo "-1"
				return
			fi
		done
		echo "0"
		return
	}
	GOVERS="${GOVERS:2}"
	EQRES=$(vercomp "$GOVERS" "1.5")
	if [ "$EQRES" == "-1" ]; then
		  echo "error: Go '1.5' or greater is required and '$GOVERS' is currently installed. Please upgrade Go at https://golang.org/dl to continue."
		  exit 1
	fi
fi

export GO15VENDOREXPERIMENT=1

cd $(dirname "${BASH_SOURCE[0]}")
OD="$(pwd)"


if [ "$1" == "vendor" ]; then
	echo "222"
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

# build and store objects into original directory.
go build -ldflags "$LDFLAGS" -o "$OD/juno-server" main/juno-server/*.go


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
	$OD/juno-server -p 6382  &
	PID=$!
	function testend {
		kill $PID &
	}
	trap testend EXIT
	go test -cover $(go list ./... | grep -v /vendor/)
fi
