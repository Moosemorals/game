
EXECUTABLE := game

SOURCEDIR := ./
BUILDDIR := ./target
BINDIR := $(BUILDDIR)

ifeq ($(OS), Windows_NT)
	MKDIR = mkdir.exe
else
	MKDIR = mkdir
endif

all: test build run

test:
	go test -v $(SOURCEDIR)

build: test
	$(MKDIR) -p $(BINDIR) 
	go build -v -o $(BINDIR)/$(EXECUTABLE) $(SOURCEDIR)

run: build
	$(BINDIR)/$(EXECUTABLE) 

clean:
	go clean -v
	rm -rf $(BUILDDIR) 

