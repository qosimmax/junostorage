all:
	@./build.sh
clean:
	rm -f juno-server
test:
	@./build.sh test
cover:
	@./build.sh cover
install: all
	sudo cp juno-server /usr/local/bin

uninstall:
	rm -f /usr/local/bin/juno-server
package:
	@./build.sh package
