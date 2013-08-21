# Requires a go dev setup with GOPATH set and code.google.com/p/tcgl/redis
# installed (go get it) in the GOPATH area.

# process to watch
PROCESS ?= atd

# override GOARCH as appropriate
GOARCH ?= 386

DEBARCH := amd64
ifeq ($(GOARCH), 386)
	DEBARCH := i386
endif

# change in ./debian/DEBIAN/control as well
VERSION := 1.0-1

all: fd-count strip

fd-count: main.go
	sed -i -e 's/~PROCESS~/${PROCESS}/' main.go
	env GOARCH=${GOARCH} go build -o fd-count-${PROCESS}
	sed -i -e 's/${PROCESS}/~PROCESS~/' main.go

strip:
	strip fd-count-${PROCESS}

deb-386:
	GOARCH=386 make deb

deb-amd64:
	GOARCH=amd64 make deb

deb: all
	sed -i -e 's/~PROCESS~/${PROCESS}/' ./debian/DEBIAN/control
	sed -i -e 's/~VERSION~/${VERSION}/' ./debian/DEBIAN/control
	sed -i -e 's/~ARCH~/${DEBARCH}/' ./debian/DEBIAN/control
	mkdir -p ./debian/usr/bin
	cp fd-count-${PROCESS} ./debian/usr/bin/
	sudo chown root:root ./debian/usr/bin/fd-count-${PROCESS}
	sudo chmod u+s ./debian/usr/bin/fd-count-${PROCESS}
	dpkg-deb --build debian
	mv debian.deb fd-count-${PROCESS}-${VERSION}_${DEBARCH}.deb
	sed -i -e 's/${DEBARCH}/~ARCH~/' ./debian/DEBIAN/control
	sed -i -e 's/${VERSION}/~VERSION~/' ./debian/DEBIAN/control
	sed -i -e 's/${PROCESS}/~PROCESS~/' ./debian/DEBIAN/control

clean-deb:
	rm -f ./debian/usr/bin/fd-count-*
	test ! -d ./debian/usr/bin || { cd ./debian; rmdir -p usr/bin; }
	rm -f fd-count-*.deb

clean: clean-deb
	rm -f fd-count-*

