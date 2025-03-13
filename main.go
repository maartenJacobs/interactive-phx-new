package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
)

var (
	projectName string
	features    []string
	database    string
	adapter     string
)

func buildCommand(
	projectName string,
	features []string,
	database string,
	adapter string,
) []string {
	featuresMap := map[string]bool{
		"Binary ID":      false,
		"Ecto":           false,
		"LiveView":       false,
		"HTML Views":     false,
		"ESBuild":        false,
		"Tailwind":       false,
		"Gettext":        false,
		"Live Dashboard": false,
		"Swoosh Mailer":  false,
	}
	for i := 0; i < len(features); i += 1 {
		featuresMap[features[i]] = true
	}

	flags := []string{
		"--install",
		"--database",
		database,
		"--adapter",
		adapter,
	}

	for feature, flag := range featuresMap {
		switch feature {
		case "Binary ID":
			if flag {
				flags = append(flags, "--binary-id")
			}
		case "Ecto":
			if !flag {
				flags = append(flags, "--no-ecto")
			}
		case "LiveView":
			if !flag {
				flags = append(flags, "--no-live")
			}
		case "HTML Views":
			if !flag {
				flags = append(flags, "--no-html")
			}
		case "ESBuild":
			if !flag {
				flags = append(flags, "--no-esbuild")
			}
		case "Tailwind":
			if !flag {
				flags = append(flags, "--no-tailwind")
			}
		case "Gettext":
			if !flag {
				flags = append(flags, "--no-gettext")
			}
		case "Live Dashboard":
			if !flag {
				flags = append(flags, "--no-dashboard")
			}
		case "Swoosh Mailer":
			if !flag {
				flags = append(flags, "--no-mailer")
			}
		}
	}

	return append(append([]string{"phx.new"}, flags...), projectName)
}

func hasMixPhxNewInstalled() (bool, error) {
	mixPath, err := exec.LookPath("mix")
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(mixPath, "phx.new")
	if err := cmd.Run(); err != nil {
		if exitError := err.(*exec.ExitError); exitError != nil {
			if exitError.ExitCode() == 1 {
				return false, nil
			}
		}
		return false, err
	}

	return true, nil
}

func ensureMixPhxNewInstalled() error {
	installed, err := hasMixPhxNewInstalled()
	if err != nil {
		return err
	}

	mixPath, err := exec.LookPath("mix")
	if err != nil {
		log.Fatal(err)
	}

	if !installed {
		fmt.Println("Phoenix not installed, running `mix archive.install phx_new`")
		return exec.Command(mixPath, "archive.install", "--force", "hex", "phx_new").Run()
	}
	return nil
}

func main() {
	if err := ensureMixPhxNewInstalled(); err != nil {
		log.Fatal(err)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What's the name of your project?").
				Prompt("? ").
				Value(&projectName),
			huh.NewMultiSelect[string]().
				Options(
					huh.NewOption("Binary ID", "Binary ID").Selected(true),
					huh.NewOption("Ecto", "Ecto").Selected(true),
					huh.NewOption("LiveView", "LiveView").Selected(true),
					huh.NewOption("HTML Views", "HTML Views").Selected(true),
					huh.NewOption("ESBuild", "ESBuild").Selected(true),
					huh.NewOption("Tailwind", "Tailwind").Selected(true),
					huh.NewOption("Gettext", "Gettext").Selected(false),
					huh.NewOption("Live Dashboard", "Live Dashboard").Selected(true),
					huh.NewOption("Swoosh Mailer", "Swoosh Mailer").Selected(false),
				).
				Title("Include Features").
				Value(&features),
			huh.NewSelect[string]().
				Title("Pick a database.").
				Options(
					huh.NewOption("Postgres", "postgres"),
					huh.NewOption("MySQL", "mysql"),
					huh.NewOption("MSSQL", "mssql"),
					huh.NewOption("SQLite3", "sqlite3"),
				).
				Value(&database),
			huh.NewSelect[string]().
				Title("Pick an adapter.").
				Options(
					huh.NewOption("Bandit", "bandit"),
					huh.NewOption("Cowboy", "cowboy"),
				).
				Value(&adapter),
		),
	)

	accessibleMode := os.Getenv("ACCESSIBLE") != ""
	form.WithAccessible(accessibleMode)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Executing the following command...")
	mixCommandArgs := buildCommand(projectName, features, database, adapter)
	fmt.Printf("mix %s\n", strings.Join(mixCommandArgs, " "))

	mixPath, err := exec.LookPath("mix")
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(mixPath, mixCommandArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
