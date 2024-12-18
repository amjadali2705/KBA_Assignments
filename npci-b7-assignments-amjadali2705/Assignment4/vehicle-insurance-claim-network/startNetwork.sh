#!/bin/sh

echo "Start  the Network"
minifab netup -s couchdb -e true -i 2.4.8 -o police.insuranceclaim.com

sleep 5

echo "create the channel"
minifab create -c insuranceclaimchannel

sleep 2

echo "Join the peers to the channel"
minifab join -c insuranceclaimchannel

sleep 2

echo "Anchor update"
minifab anchorupdate

sleep 2

echo "Profile Generation"
minifab profilegen