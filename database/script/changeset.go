package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	ROOT_PATH        = "./database/migrations/changelog"
	TEMPLATES_DIR    = "./database/templates"
	DEFAULT_TEMPLATE = "changeset-default.yml.tmpl" // ← change if needed
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run changeset.go <feature-name> <order> <changeset-name> <template-name> [var1=value1] [var2=value2] ...")
		fmt.Println("\nExample:")
		fmt.Println("  go run changeset.go RFC-2-Accounts 1 currencies changeset-default.yml.tmpl tableName=currencies")
		os.Exit(1)
	}

	featureInput := os.Args[1]
	orderStr := os.Args[2]
	changesetName := sanitizeForPath(os.Args[3])
	templateName := os.Args[4]

	// ── Parse additional key=value arguments ─────────────────────────────
	customVars := make(map[string]string)
	for i := 5; i < len(os.Args); i++ {
		parts := strings.SplitN(os.Args[i], "=", 2)
		if len(parts) == 2 {
			customVars[parts[0]] = parts[1]
		} else {
			fmt.Printf("Warning: Invalid variable format '%s', expected key=value\n", os.Args[i])
		}
	}

	// ── Resolve feature folder ──────────────────────────────────────────
	dirPath, featureFolderName, isNew, err := resolveFeatureFolder(ROOT_PATH, featureInput)
	if err != nil {
		fmt.Printf("Error resolving feature folder: %v\n", err)
		os.Exit(1)
	}

	status := "Using existing"
	if isNew {
		status = "Creating new"
	}
	fmt.Printf("→ %s feature: %s\n", status, featureFolderName)

	// ── Validate order parameter ───────────────────────────────────────
	order, err := strconv.Atoi(orderStr)
	if err != nil || order <= 0 {
		fmt.Printf("Error: Order must be a positive integer, got '%s'\n", orderStr)
		os.Exit(1)
	}

	// ── Resolve template ─────────────────────────────────────────────────
	selectedTemplate := filepath.Join(TEMPLATES_DIR, templateName)
	if _, err := os.Stat(selectedTemplate); os.IsNotExist(err) {
		fmt.Printf("Error: Template not found: %s\n", selectedTemplate)
		os.Exit(1)
	}

	fmt.Printf("→ Using template: %s\n", templateName)

	// ── Read template content ───────────────────────────────────────────
	tmplContentBytes, err := os.ReadFile(selectedTemplate)
	if err != nil {
		fmt.Printf("Cannot read template file: %v\n", err)
		os.Exit(1)
	}
	tmplContent := string(tmplContentBytes)

	// ── Extract required variables from template ────────────────────────
	requiredVars := extractTemplateVariables(tmplContent)

	// ── Prepare known values ────────────────────────────────────────────
	unixTimestamp := time.Now().UTC().Unix()
	author := getGitAuthor()

	data := map[string]string{
		"_featureFolderName": featureFolderName,
		"_featureName":       featureInput,
		"_changesetName":     changesetName,
		"_order":             fmt.Sprintf("%05d", order),
		"_unixtimestamp":     strconv.FormatInt(unixTimestamp, 10),
		"_author":            author,
		"_generatedAt":       time.Now().UTC().Format("2006-01-02 15:04:05 UTC"),
	}

	// ── Merge custom variables ──────────────────────────────────────────
	for k, v := range customVars {
		data[k] = v
	}

	// ── Check for missing required variables ────────────────────────────
	var missing []string
	for _, varName := range requiredVars {
		if _, exists := data[varName]; !exists {
			missing = append(missing, varName)
		}
	}

	if len(missing) > 0 {
		fmt.Printf("Error: Missing required variables: %s\n", strings.Join(missing, ", "))
		fmt.Println("Required variables from template:")
		for _, v := range requiredVars {
			fmt.Printf("  - %s\n", v)
		}
		os.Exit(1)
	}

	// ── Create target filename and path ─────────────────────────────────
	filename := fmt.Sprintf("%05d__%s.yml", order, changesetName)
	fullPath := filepath.Join(dirPath, filename)

	// Duplicate check
	if existing, _ := findChangesetWithSameName(dirPath, changesetName); existing != "" {
		fmt.Printf("\nERROR: Changeset with this name already exists:\n→ %s\n", filepath.Base(existing))
		os.Exit(1)
	}

	// Check for order collision
	if existingOrder, _ := findChangesetWithSameOrder(dirPath, order); existingOrder != "" {
		fmt.Printf("\nERROR: Changeset with order %05d already exists:\n→ %s\n", order, filepath.Base(existingOrder))
		os.Exit(1)
	}

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		fmt.Printf("Cannot create directory: %v\n", err)
		os.Exit(1)
	}

	// ── Parse and render template ───────────────────────────────────────
	tmplContent = strings.ReplaceAll(tmplContent, "${", "{{.")
	tmplContent = strings.ReplaceAll(tmplContent, "}", "}}")

	tmpl, err := template.New("changeset").
		Parse(tmplContent)
	if err != nil {
		fmt.Printf("Template parse error: %v\n", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		fmt.Printf("Template rendering error: %v\n", err)
		os.Exit(1)
	}

	// ── Write the final file ────────────────────────────────────────────
	if err := os.WriteFile(fullPath, buf.Bytes(), 0644); err != nil {
		fmt.Printf("Cannot write file: %v\n", err)
		os.Exit(1)
	}

	// ── Success message ─────────────────────────────────────────────────
	fmt.Println("\nChangeset created successfully!")
	fmt.Printf("  File:     %s\n", fullPath)
	fmt.Printf("  Feature:  %s\n", featureFolderName)
	fmt.Printf("  Template: %s\n", filepath.Base(selectedTemplate))
	fmt.Printf("  Author:   %s\n", author)
}

// ──────────────────────────────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────────────────────────────

func extractTemplateVariables(content string) []string {
	re := regexp.MustCompile(`\$\{([a-zA-Z0-9_]+)\}`)
	matches := re.FindAllStringSubmatch(content, -1)

	seen := make(map[string]bool)
	var vars []string

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		name := match[1]
		if !seen[name] {
			seen[name] = true
			vars = append(vars, name)
		}
	}

	sort.Strings(vars)
	return vars
}

func listAvailableTemplates() {
	files, err := os.ReadDir(TEMPLATES_DIR)
	if err != nil {
		fmt.Printf("Cannot read templates directory: %v\n", err)
		return
	}

	fmt.Println("Available templates:")
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".tmpl") {
			fmt.Printf("  - %s\n", f.Name())
		}
	}
}

// ──────────────────────────────────────────────────────────────────────
// Original helper functions (unchanged)
// ──────────────────────────────────────────────────────────────────────

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

	nextNum := maxNum + 1
	newFolderName := fmt.Sprintf("%03d_%s", nextNum, featureName)
	newPath := filepath.Join(root, newFolderName)

	return newPath, newFolderName, true, nil
}

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

func findChangesetWithSameOrder(dir string, order int) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	prefix := fmt.Sprintf("%05d__", order)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, ".yml") {
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
	re := regexp.MustCompile(`_+`)
	result = re.ReplaceAllString(result, "_")
	result = strings.Trim(result, "_")

	if result == "" {
		return "unnamed"
	}
	return result
}
