VER = 0.1.5
REV = 2
TARGET = secprac-client

$(TARGET): 
	go build -o $(TARGET)

clean:
	go clean
	rm -f $(TARGET)-$(VER)-$(REV).tar.gz

dist: $(TARGET)
	tar -z -c -f $(TARGET)-$(VER)-$(REV).tar.gz data/* $(TARGET)

.PHONY: clean dist
all: clean $(TARGET) dist
