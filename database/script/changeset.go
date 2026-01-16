package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const ROOT_PATH = "./database/migrations/changelog"

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run changeset.go <feature-name> <changeset-name>")
		os.Exit(1)
	}

	featureInput := os.Args[1]
	changesetName := sanitizeForPath(os.Args[2])
	reader := bufio.NewReader(os.Stdin)

	// Resolve feature folder (find existing or propose new)
	dirPath, featureFolderName, isNew, err := resolveFeatureFolder(ROOT_PATH, featureInput)
	if err != nil {
		fmt.Printf("Error resolving feature folder: %v\n", err)
		os.Exit(1)
	}

	// If feature already exists → ask for confirmation
	if !isNew {
		fmt.Printf("\nFound existing feature folder: %s\n", featureFolderName)
		fmt.Printf("Do you want to add the new changeset into this existing feature? [y/N]: ")

		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))

		if answer != "y" && answer != "yes" {
			fmt.Println("Aborted. Please try again with a different feature name.")
			os.Exit(0)
		}

		fmt.Printf("→ Using existing feature: %s\n", featureFolderName)
	} else {
		fmt.Printf("→ Will create new feature folder: %s\n", featureFolderName)
	}

	// Unix timestamp (seconds since 1970-01-01 UTC)
	unixTimestamp := time.Now().UTC().Unix()

	// Proposed filename
	filename := fmt.Sprintf("%d__%s.yml", unixTimestamp, changesetName)
	fullPath := filepath.Join(dirPath, filename)

	// Check if a changeset with the same logical name already exists
	existingFile, err := findChangesetWithSameName(dirPath, changesetName)
	if err != nil {
		fmt.Printf("Warning: Could not check for duplicate changeset names: %v\n", err)
	} else if existingFile != "" {
		fmt.Println("\nERROR: A changeset with this name already exists in the feature!")
		fmt.Printf("→ Found: %s\n", filepath.Base(existingFile))
		fmt.Println("Aborted. Please choose a different changeset name.")
		os.Exit(1)
	}

	// Create directory structure if it doesn't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", dirPath, err)
		os.Exit(1)
	}

	// Try to get git author name
	author := getGitAuthor()

	// Liquibase changeset template
	content := fmt.Sprintf(`databaseChangeLog:
- changeSet:
    id: %s-%d
    author: %s
    changes:
      # ──────────────────────────────────────────────────────────────
      # ↓ Put your changes here ↓
      # ──────────────────────────────────────────────────────────────
      # - createTable:
      #     tableName: example
      #     columns:
      #       - column:
      #           name: id
      #           type: bigint
      #           autoIncrement: true
      #           constraints:
      #             primaryKey: true

`, changesetName, unixTimestamp, author)

	// Write the file
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		fmt.Printf("Error writing file %s: %v\n", fullPath, err)
		os.Exit(1)
	}

	// Success output
	fmt.Println("\nCreated successfully!")
	fmt.Printf("→ %s\n", fullPath)
	fmt.Printf("  Feature: %s\n", featureFolderName)
	fmt.Printf("  Author:  %s\n", author)
}

// resolveFeatureFolder tries to find existing feature folder by name (case-insensitive)
// returns: full path, folder name, isNew, error
func resolveFeatureFolder(root, featureInput string) (string, string, bool, error) {
	featureName := sanitizeForPath(featureInput)

	files, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(root, 0755)
		} else {
			return "", "", false, err
		}
	}

	var found string
	var maxNum = 0

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		name := file.Name()
		if !strings.Contains(name, "_") {
			continue
		}

		parts := strings.SplitN(name, "_", 2)
		if len(parts) != 2 {
			continue
		}

		existingName := sanitizeForPath(parts[1])
		if strings.EqualFold(existingName, featureName) {
			found = filepath.Join(root, name)
			break
		}

		num, _ := strconv.Atoi(parts[0])
		if num > maxNum {
			maxNum = num
		}
	}

	if found != "" {
		return found, filepath.Base(found), false, nil
	}

	// Create new feature folder
	nextNum := maxNum + 1
	newFolderName := fmt.Sprintf("%03d_%s", nextNum, featureName)
	newPath := filepath.Join(root, newFolderName)

	return newPath, newFolderName, true, nil
}

// findChangesetWithSameName looks for any file ending with __<changesetName>.yml
func findChangesetWithSameName(dir, changesetName string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	suffix := "__" + changesetName + ".yml"

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, suffix) {
			return filepath.Join(dir, name), nil
		}
	}

	return "", nil
}

func getGitAuthor() string {
	out, err := exec.Command("git", "config", "--get", "user.name").Output()
	if err == nil && len(out) > 0 {
		return strings.TrimSpace(string(out))
	}

	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}
	if user == "" {
		return "unknown"
	}
	return user
}

func sanitizeForPath(s string) string {
	s = strings.TrimSpace(s)
	var sb strings.Builder

	for _, r := range s {
		switch r {
		case ' ', '/', '\\', ':', '*', '?', '"', '<', '>', '|', '\n', '\r', '\t':
			sb.WriteRune('_')
		default:
			if r >= 32 && r != 127 {
				sb.WriteRune(r)
			}
		}
	}

	result := sb.String()

	// Collapse multiple underscores and trim
	re := regexp.MustCompile(`_+`)
	result = re.ReplaceAllString(result, "_")
	result = strings.Trim(result, "_")

	if result == "" {
		return "unnamed"
	}
	return result
}