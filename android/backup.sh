#!/bin/bash
rsync --archive --compress --delete --exclude-from="$HOME/backup/backup_exclude.txt" --delete-excluded --verbose --mkpath "$HOME/storage/dcim/" backup@server:./PhoneBackup/DCIM
rsync --archive --compress --delete --exclude-from="$HOME/backup/backup_exclude.txt" --delete-excluded --verbose --mkpath "$HOME/storage/pictures/" backup@server:./PhoneBackup/Pictures
rsync --archive --compress --delete --exclude-from="$HOME/backup/backup_exclude.txt" --delete-excluded --verbose --mkpath "$HOME/storage/downloads/" backup@server:./PhoneBackup/Downloads
ssh backup@server 'cd PhoneBackup/ && git add . && git commit -m "-"'
date >> $HOME/backup/backup.log
