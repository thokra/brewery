#!/bin/sh

parent_path=$(dirname "${BASH_SOURCE}")
DEFAULT_TEMPLATE="* %h %s (%an)"

GITTEMPLATE=${GITTEMPLATE:=$DEFAULT_TEMPLATE}

if [ "$PRINT_HEADER" != "false" ]; then
	echo "# Changelog" >> /out/changelog.md
fi

git log --pretty=format:"$GITTEMPLATE" >> /out/changelog.md