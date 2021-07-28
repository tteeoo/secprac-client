#!/bin/sh

set -e

VER="0.2.1-1"

if [ "$(id -u)" -ne 0 ] ; then
	echo "run as root, e.g. 'sudo $0'"
	exit 1
fi

echo "installing secprac version $VER"

printf "creating directories... "
mkdir -p \
	/usr/local/bin \
	/var/log/secprac \
	/usr/local/share/secprac \
	/usr/local/share/fonts \
	/tmp/secprac

cd /tmp/secprac
echo "ok"

printf "downloading archive... "
if which curl > /dev/null 2>&1; then
	curl -sfLO https://github.com/tteeoo/secprac-client/releases/download/"$VER"/secprac-client-"$VER".tar.gz
else
	if which wget > /dev/null 2>&1; then
		wget https://github.com/tteeoo/secprac-client/releases/download/"$VER"/secprac-client-"$VER".tar.gz
	else
		echo "error"
		echo "either curl or wget must be installed to download the files"
	fi
fi
echo "ok"

printf "extracting archive... "
tar -z -x -f secprac-client-"$VER".tar.gz
echo "ok"

printf "installing files... "
chmod +x secprac-client data/secprac-start data/secprac-open data/secprac-report.desktop
mv -f data/*.service /etc/systemd/system/
mv -f data/*.png data/*.html /usr/local/share/secprac/
mv -f data/*.ttf /usr/local/share/fonts/
fc-cache
mv -f secprac-client data/secprac-start data/secprac-open /usr/local/bin/
if [ -n "$SUDO_USER" ] ; then
	if which xdg-user-dir > /dev/null 2>&1; then
		FILE=$(su -c 'xdg-user-dir DESKTOP' "$SUDO_USER")/secprac-report.desktop
		su -c "cp -f /tmp/secprac/data/secprac-report.desktop $FILE" "$SUDO_USER"
	else
		su -c 'cp -f /tmp/secprac/data/secprac-report.desktop $HOME/Desktop' "$SUDO_USER"
	fi
fi
echo "ok"

printf "cleaning up... "
rm -rf /tmp/secprac
echo "ok"

echo "installation successful"
