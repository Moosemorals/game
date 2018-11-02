
EXECUTABLE := game

SOURCEDIR := ./
BUILDDIR := ./target
BINDIR := $(BUILDDIR)

all: test build run

test:
	go test -v $(SOURCEDIR)

build: test
	mkdir -p $(BINDIR) 
	go build -v -o $(BINDIR)/$(EXECUTABLE) $(SOURCEDIR)

run: build
	$(BINDIR)/$(EXECUTABLE) 

clean:
	go clean -v
	rm -rf $(BUILDDIR) 

