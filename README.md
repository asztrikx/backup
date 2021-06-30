# Brief

This projects achieves backup to a server from android and linux using rsync and git for version control
Building from source requires go. If you want to avoid go you probably can easily rewrite it in bash if you have experience with it.

# Install

Make copies of _\*\_template.\*_ files without this suffix

## PC

I will run this as root to backup all users' files. If you want to run as normal user then change paths, systemd service, ...

-   Extend /etc/ssh/config: `sudo nano /etc/ssh/config` //make address aliases accessible for everyone (e.g. root)
    ```
    Host server
        Hostname ...
    ```
-   Build main.go
    -   `cd pc && go build -o backup && cd ..`
    -   `chmod u+x pc/backup` //make it executable
-   Move files
    -   `sudo mv pc/backup /usr/bin/`
    -   `sudo cp pc/backup_translation.json /etc/`
    -   `sudo cp pc/backup_exclude.txt /etc/`
    -   `sudo cp pc/backup.service /etc/systemd/system/`
    -   `sudo cp pc/backup.timer /etc/systemd/system/`
    -   `sudo ln --symbolic ~/.ssh/id_rsa /etc/backup_id_rsa`
-   Edit backup_translation.json to your needs: `sudo code /etc/backup_translation.json`
    -   if target is empty then the absolute path will be recreated under the backup folder specified in main.go
-   Edit backup_exclude.txt to your needs (see rsync manual): `sudo code /etc/backup_exclude.txt`
-   Edit backup.timer to your needs (see systemd timers): `sudo nano /etc/systemd/system/backup.timer`
-   Enable service
    -   `sudo systemctl enable backup.timer`
    -   `sudo systemctl start backup.timer`

## Android

-   On phone
    -   [Install fdroid](https://wiki.termux.com/wiki/Termux_Google_Play)
        -   Install termux
    -   [termux](https://wiki.termux.com/wiki)
        -   `pkg upgrade` //get latest packages
        -   `pkg install openssh` //ssh for phone for easier usage and for rsync
        -   `echo "ssh-rsa ..." >> .ssh/authorized_keys` //inside "" should be the public key (put it into [Keep](keep.google.com) for example and copy to here)
        -   `sed --in-place "s/PasswordAuthentication yes/PasswordAuthentication no/" $HOME/../usr/etc/ssh/sshd_config` //disable password login for better security
        -   `sshd` //start ssh
        -   `whoami` //remember this username for ssh config
-   On linux

    -   Extend ~/.ssh/config if it isn't already: `nano ~/.ssh/config`

        ```
        Host phone
        	User ...
        	Hostname ...
        	Port 8022

        ```

    -   Make phone ready: `ssh phone`

        -   `termux-setup-storage` //request permission to access internal storage so we can back it up
        -   `ssh-keygen -t rsa -b 4096` //create public and private key for phone to be able to connect to server
        -   `pkg install rsync` //install rsync
        -   `pkg install termux-services` //enables termux to not get killed in the background
        -   `pkg install cronie` //regularly run our script with crontab
        -   Extend ~/.ssh/config: `nano ~/.ssh/config`

            ```
            Host server
            	Hostname ...

            ```

        -   `exit`

    -   Move files to and from phone

        -   `rsync --mkpath android/backup.sh phone:backup/`
        -   `rsync --mkpath android/backup_exclude.txt phone:backup/`
        -   `rsync --mkpath android/backup_cronjob.txt phone:backup/`
        -   `rsync phone:.ssh/id_rsa.pub /tmp/id_rsa.pub`
        -   `rsync /tmp/id_rsa.pub backup@server:./` //the alias backup@server should correspond to a separate account for this on the backup server

    -   Make server ready: `ssh backup@server`

        -   `cat id_rsa.pub >> .ssh/known_hosts` //add copied public key to known hosts
        -   `rm id_rsa.pub` //we don't need it anymore
        -   `exit`

    -   Start syncing : `ssh phone`
        -   `nano backup/backup_cronjob.txt` //modify according to your needs (currently every minute the sync will run if the previous sync is finished)
        -   `nano backup/backup_exclude.txt` //modify according to your needs (check rsync manual)
        -   `sv-enable crond` //make crontab autostart
        -   `crontab backup/backup_cronjob.txt` //register our crontab config
        -   `exit`

# Todo

-   Replace bash on android if more directories have to be synced
