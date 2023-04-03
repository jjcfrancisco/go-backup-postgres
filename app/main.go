package main

import (
	"fmt"
	"os"
)

func main() {

	available := availableBackups()

	remove := checkBackups(available)

	err := makeBackup()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	destroyBackup(remove)

	fmt.Println("Backup completed successfully.")
}
