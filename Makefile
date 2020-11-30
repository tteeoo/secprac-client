VER = 0.1.5
TARGET = secprac-client

$(TARGET): 
	go build -o $(TARGET)

clean:
	go clean
	rm -f $(TARGET)-$(VER).tar.gz

dist: $(TARGET)
	tar -z -c -f $(TARGET)-$(VER).tar.gz data/* $(TARGET)

.PHONY: clean dist
all: clean $(TARGET) dist
