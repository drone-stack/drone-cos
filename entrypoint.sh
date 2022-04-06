#!/bin/sh
set -eu

[ ! -z "$1" ] && (
    exec /bin/bash
) 

[ -z "${PLUGIN_DEBUG}" ] || set -x

[ -z "${PLUGIN_PAUSE}" ] || sleep 100000

if [ ! -z "${PLUGIN_PROXY}" ]; then
	export http_proxy=${PLUGIN_PROXY}
	export https_proxy=${PLUGIN_PROXY}
	export all_proxy=${PLUGIN_PROXY}
	export no_proxy=localhost,127.0.0.1/8
	echo "http proxy done"
fi

exec "$@"