package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type BitwardenItem struct {
	ID         string      `json:"id"`
	FolderID   string      `json:"folderId"`
	Type       int         `json:"type"`
	Name       string      `json:"name"`
	Notes      string      `json:"notes"`
	Login      *Login      `json:"login"`
	SecureNote *SecureNote `json:"secureNote"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	URIs     []URI  `json:"uris"`
	TOTP     string `json:"totp"`
}

type URI struct {
	Match interface{} `json:"match"`
	URI   string      `json:"uri"`
}

type SecureNote struct {
	Type int `json:"type"`
}

type Folder struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type BitwardenExport struct {
	Items   []BitwardenItem `json:"items"`
	Folders []Folder        `json:"folders"`
}

func sanitizeName(name string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, name)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: bw2pass <bitwarden_export.json>")
		os.Exit(1)
	}

	jsonFile := os.Args[1]
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	var bwExport BitwardenExport
	err = json.Unmarshal(data, &bwExport)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	folderMap := make(map[string]string)
	for _, folder := range bwExport.Folders {
		folderMap[folder.ID] = sanitizeName(folder.Name)
	}

	entryCount := make(map[string]int)

	for _, item := range bwExport.Items {
		var passName string
		var passContent string

		sanitizedName := sanitizeName(item.Name)

		if item.Type == 1 { // Login type
			var domain string
			if item.Login != nil && len(item.Login.URIs) > 0 {
				domain = getDomainFromURI(item.Login.URIs[0].URI)
			} else {
				domain = "unknown_domain"
			}

			if item.FolderID != "" && item.FolderID != "null" {
				folderName := folderMap[item.FolderID]
				passName = fmt.Sprintf("%s/%s/%s", folderName, domain, sanitizedName)
			} else {
				passName = fmt.Sprintf("%s/%s", domain, sanitizedName)
			}

			passContent = fmt.Sprintf("%s\nUsername: %s\n", item.Login.Password, item.Login.Username)

			for _, uri := range item.Login.URIs {
				passContent += fmt.Sprintf("URL: %s\n", uri.URI)
			}

			if item.Login.TOTP != "" {
				passContent += fmt.Sprintf("TOTP: %s\n", item.Login.TOTP)
			}
		} else if item.Type == 2 { // Secure Note type
			if item.FolderID != "" && item.FolderID != "null" {
				folderName := folderMap[item.FolderID]
				passName = fmt.Sprintf("%s/notes/%s", folderName, sanitizedName)
			} else {
				passName = fmt.Sprintf("notes/%s", sanitizedName)
			}

			passContent = item.Notes
		}

		if item.Notes != "" && item.Type == 1 {
			passContent += fmt.Sprintf("\nNotes:\n%s\n", item.Notes)
		}

		// Check if this path combination already exists
		if count, exists := entryCount[passName]; exists {
			entryCount[passName] = count + 1
			passName = fmt.Sprintf("%s_%d", passName, count+1)
		} else {
			entryCount[passName] = 1
		}

		cmd := exec.Command("pass", "insert", "-m", passName)
		cmd.Stdin = strings.NewReader(passContent)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error inserting %s: %v\n", passName, err)
		} else {
			fmt.Printf("Inserted: %s\n", passName)
		}
	}
}

func getDomainFromURI(uri string) string {
	parts := strings.Split(uri, "://")
	if len(parts) > 1 {
		domain := strings.Split(parts[1], "/")[0]
		return domain
	}
	return "unknown_domain"
}
