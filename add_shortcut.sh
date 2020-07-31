#!/bin/bash
export GOPATH=/mnt/Data/SteaMyProjects/go
file=/usr/share/applications/swap.desktop
echo "[Desktop Entry]
Name=SWAutoPlay
Terminal=true
Type=Application
Icon="$GOPATH"/src/SWAutoPlay_GUI/data/icon.ico
Exec="$GOPATH"/src/SWAutoPlay_GUI/run.sh" > $file
update-desktop-database