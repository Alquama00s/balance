package main

import (
	"bufio"
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
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run changeset.go <feature-name> <changeset-name>")
		os.Exit(1)
	}

	featureInput := os.Args[1]
	changesetName := sanitizeForPath(os.Args[2])
	reader := bufio.NewReader(os.Stdin)

	// ── Resolve feature folder ──────────────────────────────────────────
	dirPath, featureFolderName, isNew, err := resolveFeatureFolder(ROOT_PATH, featureInput)
	if err != nil {
		fmt.Printf("Error resolving feature folder: %v\n", err)
		os.Exit(1)
	}

	if !isNew {
		fmt.Printf("\nFound existing feature: %s\n", featureFolderName)
		fmt.Print("Add to this feature? [Y/n]: ")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer == "n" || answer == "no" {
			fmt.Println("Aborted.")
			os.Exit(0)
		}
		fmt.Printf("→ Using existing feature: %s\n", featureFolderName)
	} else {
		fmt.Printf("→ Will create new feature folder: %s\n", featureFolderName)
	}

	// ── Select template ─────────────────────────────────────────────────
	selectedTemplate, err := selectTemplate(reader)
	if err != nil {
		fmt.Printf("Template selection failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("→ Using template: %s\n", filepath.Base(selectedTemplate))

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
		"_unixtimestamp":     strconv.FormatInt(unixTimestamp, 10),
		"_author":            author,
		"_generatedAt":       time.Now().UTC().Format("2006-01-02 15:04:05 UTC"),
	}
	for k, v := range data {
		println(k + ": " + v)
	}

	// ── Prompt for missing variables ────────────────────────────────────
	fmt.Println("\nPlease provide values for the following variables:")
	for _, varName := range requiredVars {
		// Skip already provided variables
		if _, exists := data[varName]; exists {
			continue
		}

		fmt.Printf("%s: ", varName)
		value, _ := reader.ReadString('\n')
		value = strings.TrimSpace(value)

		// You may want to add a default here if value == ""
		// value = value or "default_value"

		data[varName] = value
	}

	// ── Create target filename and path ─────────────────────────────────
	filename := fmt.Sprintf("%d__%s.yml", unixTimestamp, changesetName)
	fullPath := filepath.Join(dirPath, filename)

	// Duplicate check
	if existing, _ := findChangesetWithSameName(dirPath, changesetName); existing != "" {
		fmt.Printf("\nERROR: Changeset with this name already exists:\n→ %s\n", filepath.Base(existing))
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

func selectTemplate(reader *bufio.Reader) (string, error) {
	files, err := os.ReadDir(TEMPLATES_DIR)
	if err != nil {
		return "", fmt.Errorf("cannot read templates directory %s: %w", TEMPLATES_DIR, err)
	}

	var templates []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if strings.HasSuffix(name, ".tmpl") {
			templates = append(templates, name)
		}
	}

	if len(templates) == 0 {
		return "", fmt.Errorf("no .tmpl files found in %s", TEMPLATES_DIR)
	}

	sort.Strings(templates)

	// Find default index
	defaultIdx := 0
	for i, name := range templates {
		if name == DEFAULT_TEMPLATE {
			defaultIdx = i
			break
		}
	}

	fmt.Println("\nAvailable templates:")
	for i, name := range templates {
		prefix := "  "
		if i == defaultIdx {
			prefix = "* "
		}
		fmt.Printf("%s%d) %s\n", prefix, i+1, name)
	}

	fmt.Printf("\nSelect template (1–%d) [default = %d (%s)]: ",
		len(templates), defaultIdx+1, templates[defaultIdx])

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return filepath.Join(TEMPLATES_DIR, templates[defaultIdx]), nil
	}

	num, err := strconv.Atoi(input)
	if err != nil || num < 1 || num > len(templates) {
		fmt.Printf("Invalid selection → using default (%s)\n", templates[defaultIdx])
		return filepath.Join(TEMPLATES_DIR, templates[defaultIdx]), nil
	}

	return filepath.Join(TEMPLATES_DIR, templates[num-1]), nil
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
