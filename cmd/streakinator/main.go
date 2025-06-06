// main.go
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(`
  _________ __                         __   .__               __                
 /   _____//  |________   ____ _____  |  | _|__| ____ _____ _/  |_  ___________ 
 \_____  \\   __\_  __ \_/ __ \\__  \ |  |/ /  |/    \\__  \\   __\/  _ \_  __ \
 /        \|  |  |  | \/\  ___/ / __ \|    <|  |   |  \/ __ \|  | (  <_> )  | \/
/_______  /|__|  |__|    \___  >____  /__|_ \__|___|  (____  /__|  \____/|__|   
        \/                   \/     \/     \/       \/     \/                   
	`)
	fmt.Println("This tool will create backdated, empty commits to populate a GitHub contribution graph.")
	fmt.Println()
	fmt.Println("WARNING:")
	fmt.Println("This utility is provided strictly for emergency or maintenance purposes—such as recovering lost commit history, migrating an existing codebase, or filling genuine gaps in a personal project’s timeline. Under NO CIRCUMSTANCES should you use this tool to fabricate GitHub activity for résumé inflation. Misusing this tool to simulate work you did not actually complete is unethical and could damage your professional reputation or lead to legal/academic consequences.")
	fmt.Println()
	fmt.Print("Press ENTER to acknowledge and continue, or CTRL+C to abort: ")
	reader.ReadString('\n')

	// STEP 1: Repository Directory
	fmt.Print("Enter the path to your Git repository (leave blank for current directory): ")
	dirInput, _ := reader.ReadString('\n')
	repoPath := strings.TrimSpace(dirInput)
	if repoPath == "" {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error: Unable to determine current directory:", err)
			return
		}
		repoPath = pwd
	}

	// If path does not exist, create it
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		fmt.Printf("The path '%s' does not exist. Creating directory.\n", repoPath)
		if err := os.MkdirAll(repoPath, 0755); err != nil {
			fmt.Println("Error: Unable to create directory:", err)
			return
		}
	}

	// Change working directory
	if err := os.Chdir(repoPath); err != nil {
		fmt.Println("Error: Unable to change directory to", repoPath, ":", err)
		return
	}
	fmt.Println("Working directory set to:", repoPath)

	// If there is no .git folder, initialize a new repository
	gitDir := filepath.Join(repoPath, ".git")
	if info, err := os.Stat(gitDir); os.IsNotExist(err) || !info.IsDir() {
		fmt.Println("No Git repository detected. Initializing with 'git init'.")
		if err := runGit("init"); err != nil {
			fmt.Println("Error: 'git init' failed:", err)
			return
		}
		fmt.Println("Git repository successfully initialized.")
	}

	// STEP 2: Branch Selection
	fmt.Print("\nEnter the branch name to which commits will be added: ")
	branchInput, _ := reader.ReadString('\n')
	branch := strings.TrimSpace(branchInput)
	if branch == "" {
		fmt.Println("Error: Branch name cannot be blank.")
		return
	}

	if !branchExists(branch) {
		fmt.Printf("Branch '%s' does not exist.\n", branch)
		fmt.Print("Would you like to create it? (y/n): ")
		resp, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(resp)) != "y" {
			fmt.Println("Operation aborted. A valid branch is required.")
			return
		}
		if err := runGit("checkout", "-b", branch); err != nil {
			fmt.Println("Error: Unable to create branch:", err)
			return
		}
		fmt.Println("Branch created and checked out:", branch)
	} else {
		fmt.Println("Checking out existing branch:", branch)
		if err := runGit("checkout", branch); err != nil {
			fmt.Println("Error: Unable to check out branch:", err)
			return
		}
		fmt.Println("Branch successfully checked out:", branch)
	}

	// STEP 3: Commit Mode Selection
	fmt.Println("\nSelect commit mode:")
	fmt.Println("  1) Fixed commits per day")
	fmt.Println("  2) Randomized commits per day")
	fmt.Print("Enter choice (1 or 2): ")
	modeInput, _ := reader.ReadString('\n')
	mode := strings.TrimSpace(modeInput)

	var fixedCount, minCount, maxCount int
	if mode == "1" {
		fmt.Print("Enter the number of commits to create each day: ")
		fixedStr, _ := reader.ReadString('\n')
		val, err := strconv.Atoi(strings.TrimSpace(fixedStr))
		if err != nil || val <= 0 {
			fmt.Println("Error: Invalid number. Operation aborted.")
			return
		}
		fixedCount = val
	} else if mode == "2" {
		fmt.Print("Enter the minimum number of commits per day: ")
		minStr, _ := reader.ReadString('\n')
		minVal, err := strconv.Atoi(strings.TrimSpace(minStr))
		if err != nil || minVal < 0 {
			fmt.Println("Error: Invalid minimum. Operation aborted.")
			return
		}

		fmt.Print("Enter the maximum number of commits per day: ")
		maxStr, _ := reader.ReadString('\n')
		maxVal, err := strconv.Atoi(strings.TrimSpace(maxStr))
		if err != nil || maxVal < minVal {
			fmt.Println("Error: Invalid maximum (must be greater than or equal to minimum). Operation aborted.")
			return
		}

		minCount = minVal
		maxCount = maxVal
	} else {
		fmt.Println("Error: Invalid selection. Operation aborted.")
		return
	}

	// STEP 4: Date Selection
	fmt.Println("\nSelect date configuration:")
	fmt.Println("  1) Single date")
	fmt.Println("  2) Date range")
	fmt.Print("Enter choice (1 or 2): ")
	dateChoice, _ := reader.ReadString('\n')
	dateChoice = strings.TrimSpace(dateChoice)

	var timestamps []time.Time
	if dateChoice == "1" {
		fmt.Print("Enter the date (YYYY-MM-DD): ")
		dateStr, _ := reader.ReadString('\n')
		dateStr = strings.TrimSpace(dateStr)
		day, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
		if err != nil {
			fmt.Println("Error: Invalid date format. Operation aborted.")
			return
		}

		var count int
		if mode == "1" {
			count = fixedCount
		} else {
			count = rand.Intn(maxCount-minCount+1) + minCount
		}

		for i := 0; i < count; i++ {
			ts := time.Date(day.Year(), day.Month(), day.Day(), 12, 0, 0, 0, time.Local).Add(time.Duration(i) * time.Second)
			timestamps = append(timestamps, ts)
		}
	} else if dateChoice == "2" {
		fmt.Print("Enter the start date (YYYY-MM-DD): ")
		startStr, _ := reader.ReadString('\n')
		startStr = strings.TrimSpace(startStr)
		startDate, err := time.ParseInLocation("2006-01-02", startStr, time.Local)
		if err != nil {
			fmt.Println("Error: Invalid start date. Operation aborted.")
			return
		}

		fmt.Print("Enter the end date (YYYY-MM-DD): ")
		endStr, _ := reader.ReadString('\n')
		endStr = strings.TrimSpace(endStr)
		endDate, err := time.ParseInLocation("2006-01-02", endStr, time.Local)
		if err != nil {
			fmt.Println("Error: Invalid end date. Operation aborted.")
			return
		}
		if endDate.Before(startDate) {
			fmt.Println("Error: End date is before start date. Operation aborted.")
			return
		}

		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			var count int
			if mode == "1" {
				count = fixedCount
			} else {
				count = rand.Intn(maxCount-minCount+1) + minCount
			}
			for i := 0; i < count; i++ {
				ts := time.Date(d.Year(), d.Month(), d.Day(), 12, 0, 0, 0, time.Local).Add(time.Duration(i) * time.Second)
				timestamps = append(timestamps, ts)
			}
		}
	} else {
		fmt.Println("Error: Invalid selection. Operation aborted.")
		return
	}

	// STEP 5: Commit Message
	fmt.Print("\nEnter commit message (press Enter for default): ")
	msgInput, _ := reader.ReadString('\n')
	commitMsg := strings.TrimSpace(msgInput)
	if commitMsg == "" {
		commitMsg = "Streakinator commit"
	}

	// STEP 6: Create Backdated Commits
	fmt.Println("\nCreating backdated commits...")
	for i, ts := range timestamps {
		dateISO := ts.Format("2006-01-02T15:04:05-07:00")
		if err := createCommit(dateISO, commitMsg); err != nil {
			fmt.Printf("Error: Commit %d/%d at %s failed: %v\n", i+1, len(timestamps), dateISO, err)
			return
		}
		fmt.Printf("Commit %d/%d created at %s\n", i+1, len(timestamps), dateISO)
	}

	fmt.Println("\nAll commits have been created successfully.")
	fmt.Printf("To upload your commits, run: git push origin %s\n", branch)
}

// runGit executes a git command and prints its output.
func runGit(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// branchExists returns true if the specified branch exists locally.
func branchExists(branch string) bool {
	cmd := exec.Command("git", "rev-parse", "--verify", branch)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}

// createCommit sets GIT_AUTHOR_DATE and GIT_COMMITTER_DATE, then makes an empty commit.
func createCommit(dateISO, message string) error {
	cmd := exec.Command("git", "commit", "--allow-empty", "-m", message)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_DATE="+dateISO,
		"GIT_COMMITTER_DATE="+dateISO,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
