#!/bin/sh
mkdir -p /usr/local/bin
mkdir -p /var/log/secprac
mkdir -p /usr/local/share/secprac
mkdir -p /usr/local/lib/systemd/system
cd /usr/local/share/secprac
rm secprac-plus.png secprac-minus.png secprac-info.png
wget https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-plus.png
wget https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-minus.png
wget https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-info.png
cd /usr/local/bin
rm secprac-client
wget https://github.com/blueberry-jam/secprac-client/releases/download/0.1.4/secprac-client
chmod +x secprac-client
