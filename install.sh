#!/bin/sh

set -e

VER="0.1.5-2"

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
	/tmp/secprac

cd /tmp/secprac
echo "ok"

printf "downloading archive... "
curl -sfLO https://github.com/blueberry-jam/secprac-client/releases/download/"$VER"/secprac-client-"$VER".tar.gz
echo "ok"

printf "extracting archive... "
tar -z -x -f secprac-client-"$VER".tar.gz
echo "ok"

printf "installing files... "
chmod +x secprac-client data/secprac-start
mv -f data/*.service /etc/systemd/system/
mv -f data/*.png /usr/local/share/secprac/
mv -f secprac-client data/secprac-start /usr/local/bin/
echo "ok"

printf "cleaning up... "
rm -rf /tmp/secprac
echo "ok"

echo "installation successful"
