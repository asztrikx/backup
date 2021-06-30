package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/asztrikx/go-terminal"
)

type translation struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

//readJSON reads a json file to `v` so it must be passed as reference
func readJSON(path string, v interface{}) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	if err := json.NewDecoder(file).Decode(v); err != nil {
		log.Fatalln(err)
	}
	if err := file.Close(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	//read config to map source path to target path
	var translationS []translation
	readJSON("/etc/backup_translation.json", &translationS)

	//send files to server with rsync
	for i, translation := range translationS {
		source := translation.Source
		target := translation.Target

		//if target is empty then copy source as target
		//trailing / would cause duplicate folder so remove it (if exists), check rsync manual for further information
		if target == "" {
			target = strings.TrimSuffix(source, "/")
		}

		//based on previous check on trailing /: disable this, otherwise correct syntax, to avoid falling into this mistake
		//if source is file (has no trailing slash) then trailing slash for target is allowed
		if strings.HasSuffix(source, "/") && strings.HasSuffix(target, "/") {
			log.Printf("malformed target directory in %d. config: %s\n", i, target)
			continue
		}

		//optinal: avoid double / by removing leading / in target
		target = strings.TrimPrefix(target, "/")

		//execute command
		//--rhs: use specified private key
		//--archive: recursion + preserve permissions
		//--compress: compress while transporting
		//--delete: remove extra files not on source
		//--exclude-from: lines of patterns to exclude
		//--delete-excluded: remove extra files even if they are excluded (one might later add files to exclude so they should be deleted from server to have 1:1 copy)
		//--verbose: for debug
		//--mkpath: create folder structure for files
		command := fmt.Sprintf(`rsync --rsh="ssh -i /etc/backup_id_rsa" --archive --compress --delete --exclude-from='/etc/backup_exclude.txt' --delete-excluded --verbose --mkpath %s backup@server:./LaptopBackup/%s`, source, target)
		_, stdErr := terminal.Bash.Exec(command)
		if stdErr != "" {
			log.Println(stdErr)
		}
	}

	//create a version of the modified files with git
	terminal.Bash.Exec(`ssh -n -i /etc/backup_id_rsa backup@server 'cd ./LaptopBackup/ && git add . && git commit -m "-"'`)
}
