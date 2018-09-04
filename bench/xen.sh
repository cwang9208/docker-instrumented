#!/bin/bash

DOMID=$( xl list | grep "^${DOMAINNAME}" | awk '{ print $2 }' )

CONDEV=$( xenstore-read /local/domain/${DOMID}/console/tty )
echo " Domain console: $CONDEV"
