package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Gets environment variables
var (
	USERNAME string = os.Getenv("USERNAME")
	DATABASE string = os.Getenv("DATABASE")
	PASSWORD string = os.Getenv("PASSWORD")
	HOST     string = os.Getenv("HOST")
	PORT     string = os.Getenv("PORT")
	TABLE    string = os.Getenv("TABLE")
)

func makeBackup() error {
	// Uses pg_dump to produce a '.dump' backup of a given table.

	now := time.Now()
	layout := "010206__1504"
	date := strings.ToUpper(now.Format(layout))
	pgDumpCmd := fmt.Sprintf("PGPASSWORD=%s pg_dump -h %s -p %s -Fc -U %s -d %s -t %s > 'backups/backup_%s.dump'", PASSWORD, HOST, PORT, USERNAME, DATABASE, TABLE, date)

	cmd := exec.Command("bash", "-c", pgDumpCmd)

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func checkBackups(b []*Backups) *Backups {
	// Checks for 'expired' backups that need to be deleted.

	tz, err := time.LoadLocation("Europe/London")
	if err != nil {
		fmt.Println(err)
	}

	now := time.Now().In(tz)
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, now.Location())

	nowMinusSix := today.AddDate(0, 0, -6)
	year, month, day = nowMinusSix.Date()
	sixDaysAgo := time.Date(year, month, day, 0, 0, 0, 0, nowMinusSix.Location())

	for _, file := range b {

		date := file.Date
		layout := "010206"

		dateParsed, err := time.Parse(layout, date)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		dateParsed = dateParsed.In(tz)
		year, month, day = dateParsed.Date()
		backupDate := time.Date(year, month, day, 0, 0, 0, 0, dateParsed.Location())

		if backupDate.Before(sixDaysAgo) {

			return file
		}
	}

	return nil

}

func destroyBackup(b *Backups) error {
	// Removes 'expired' backups

	filename := b.FileName
	err := os.Remove(filename)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func availableBackups() []*Backups {
	// Shows all current backup '.dump' files

	dir := "backups/"
	files, err := filepath.Glob(filepath.Join(dir, "*.dump"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	backups := []*Backups{}

	for _, file := range files {
		fileSplit := strings.Split(file, "/")
		fileNameWithExtension := fileSplit[len(fileSplit)-1]
		fileNameNoExtension := strings.Replace(fileNameWithExtension, ".dump", "", 1)
		timestamp := strings.Replace(fileNameNoExtension, "backup_", "", 1)
		timestampSplit := strings.Split(timestamp, "__")
		dateOnly := timestampSplit[0]
		backup := Backups{FileName: file, Date: dateOnly}
		backups = append(backups, &backup)
	}

	return backups
}
