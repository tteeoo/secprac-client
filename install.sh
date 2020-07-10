mkdir -p /usr/local/bin /var/log/secprac /etc/systemd/system /usr/local/share/secprac
cd /etc/systemd/system
rm secprac-client@.service
wget https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-client@.service
cd /usr/local/bin
rm secprac-client secprac-start
wget https://github.com/blueberry-jam/secprac-client/releases/download/0.1.5/secprac-client https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-start
chmod +x secprac-client secprac-start
cd /usr/local/share/secprac
rm secprac-plus.png secprac-minus.png secprac-info.png team
wget https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-info.png https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-plus.png https://raw.githubusercontent.com/blueberry-jam/secprac-client/master/data/secprac-minus.png
