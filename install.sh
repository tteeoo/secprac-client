#!/bin/sh
mkdir -p /usr/local/bin
mkdir -p /var/log/secprac
mkdir -p /etc/systemd/system
mkdir -p /usr/local/share/secprac
cd /usr/local/bin
rm secprac-client secprac-start
wget https://github.com/blueberry-jam/secprac-client/releases/download/0.1.4/secprac-client
wget https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-start
chmod +x secprac-client secprac-start
cd /usr/local/share/secprac
rm secprac-plus.png secprac-minus.png secprac-info.png
wget https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-info.png
wget https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-plus.png
wget https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-minus.png
cd /etc/systemd/system
rm secprac-client@.service
wget https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-client@.service
