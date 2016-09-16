all:
	@./build.sh
clean:
	rm -f server
test:
	@./build.sh test
cover:
	@./build.sh cover
install: all
	sudo cp server /usr/local/bin

uninstall:
	rm -f /usr/local/bin/server
package:
	@./build.sh package
