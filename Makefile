VER = 0.1.9
REV = 1
TARGET = secprac-client

$(TARGET): 
	CGO_ENABLED=0 go build -o $(TARGET)

clean:
	go clean
	rm -f $(TARGET)-$(VER)-$(REV).tar.gz

dist: $(TARGET)
	tar -z -c -f $(TARGET)-$(VER)-$(REV).tar.gz data/* $(TARGET)

install: $(TARGET)
	mkdir -p /usr/local/bin /var/log/secprac /usr/local/share/secprac
	touch /usr/local/share/secprac/report.html
	cp -f $(TARGET) data/secprac-start data/secprac-open /usr/local/bin/
	cp -f data/*.png data/*.html /usr/local/share/secprac/
	if which systemctl > /dev/null 2>&1; then\
		mkdir -p /etc/systemd/system;\
		cp -f data/*.service /etc/systemd/system/;\
	fi

.PHONY: clean dist
all: clean $(TARGET) dist
