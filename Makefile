VER = 0.1.5
REV = 3
TARGET = secprac-client

$(TARGET): 
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o test -ldflags '-extldflags "-f no-PIC -static"' -buildmode pie -tags 'osusergo netgo static_build' -o $(TARGET)

clean:
	go clean
	rm -f $(TARGET)-$(VER)-$(REV).tar.gz

dist: $(TARGET)
	tar -z -c -f $(TARGET)-$(VER)-$(REV).tar.gz data/* $(TARGET)

.PHONY: clean dist
all: clean $(TARGET) dist
