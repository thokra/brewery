#!/bin/sh
if [ -n "$LAGER_NETRC_MACHINE" ]; then
cat <<EOF > $HOME/.netrc
machine $LAGER_NETRC_MACHINE
login $LAGER_NETRC_USERNAME
password $LAGER_NETRC_PASSWORD
EOF
chmod 0600 $HOME/.netrc
fi
unset LAGER_NETRC_USERNAME
unset LAGER_NETRC_PASSWORD

cat $HOME/.netrc

if [ -n "$LAGER_MOUNT_DIR" ]; then
	mkdir -p "$LAGER_MOUNT_DIR"
fi

if [ -n "$LAGER_REPO" ]; then
	git clone $LAGER_REPO $LAGER_MOUNT_DIR
fi
