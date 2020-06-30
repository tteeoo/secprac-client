#!/bin/sh
mkdir -p /usr/local/bin
mkdir -p /var/log/secprac
mkdir -p /usr/local/share/secprac
mkdir -p /usr/local/lib/systemd/system
cd /usr/local/share/secprac
wget https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-plus.png
wget https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-minus.png
wget https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-info.png
wget https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-logo.png
cd /usr/local/bin
wget https://github.com/blueberry-jam/secprac-client/releases/download/0.1.1/secprac-client
chmod +x secprac-client
