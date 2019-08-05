package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	scalingo "github.com/Scalingo/go-scalingo"
)

func main() {

	// ---- PARSE ARGS ----
	if len(os.Args) != 5 {
		fmt.Fprintf(os.Stderr, "Usage: ./backup_download [REGION] [APP_NAME] [ADDON_ID] [BACKUP_ID]\n")
		os.Exit(1)
	}

	region := os.Args[1]
	appName := os.Args[2]
	addonId := os.Args[3]
	backupId := os.Args[4]

	// ---- CREATE SCALINGO CLIENT ----
	token := os.Getenv("SCALINGO_TOKEN")
	if token == "" {
		fmt.Fprintf(os.Stderr, "SCALINGO_TOKEN is not set, please set it before using this example\n")
		os.Exit(1)
	}

	client, err := scalingo.New(scalingo.ClientConfig{
		Region:   region,
		APIToken: token,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to create scalingo client: %s\n", err.Error())
		os.Exit(1)
	}

	// ---- GET BACKUP DOWNLOAD URL ----
	backupURL, err := client.BackupDownloadURL(appName, addonId, backupId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to get backup URL: %s\n", err.Error())
		os.Exit(1)
	}

	// ---- START DOWNLOAD TO ./backup.tar.gz ----

	// Create output file
	file, err := os.Create("backup.tar.gz")
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to open backup file: %s\n", err.Error())
		os.Exit(1)
	}
	defer file.Close()

	// Start HTTP request
	resp, err := http.Get(backupURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to start backup download: %s\n", err.Error())
	}
	defer resp.Body.Close()

	fmt.Println("Start backup download.")

	// Copy request body to output file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to download backup: %s\n", err.Error())
	}

	fmt.Println("Backup saved to ./archive.tar.gz")
}
