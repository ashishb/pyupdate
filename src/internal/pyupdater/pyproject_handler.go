package pyupdater

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"strings"
)

func getPyProjectTomlData(pyprojectPath string) ([]byte, error) {
	// Check if the pyproject.toml file exists
	if _, err := os.Stat(pyprojectPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("pyproject.toml not found at %s", pyprojectPath)
	}

	log.Debug().
		Str("pyprojectPath", pyprojectPath).
		Msg("Found pyproject.toml file")
	data, err := os.ReadFile(pyprojectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read pyproject.toml: %w", err)
	}

	return data, nil
}

func removeDependenciesFromTomlFile(data []byte, pyprojectPath string) error {
	var pyProjectToml map[string]any
	if err := toml.Unmarshal(data, &pyProjectToml); err != nil {
		return fmt.Errorf("failed to unmarshal pyproject.toml: %w", err)
	}

	if project, ok := pyProjectToml["project"].(map[string]any); ok {
		log.Info().Msg("Removing main dependencies from pyproject.toml")
		delete(project, "dependencies")
	}
	if depGroups, ok := pyProjectToml["dependency-groups"].(map[string]any); ok {
		log.Info().Msg("Removing dev dependencies from pyproject.toml")
		delete(depGroups, "dev")
	}

	updatedData, err := toml.Marshal(pyProjectToml)
	if err != nil {
		return fmt.Errorf("failed to marshal updated pyproject.toml: %w", err)
	}
	if err := os.WriteFile(pyprojectPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write updated pyproject.toml: %w", err)
	}
	return nil
}

func parsePyProjectToml(data []byte) (*_PyProjectToml, error) {
	var pyproject _PyProjectToml
	if err := toml.Unmarshal(data, &pyproject); err != nil {
		return nil, fmt.Errorf("failed to parse pyproject.toml: %w", err)
	}
	return &pyproject, nil
}

func withoutVersion(deps []string) []string {
	result := make([]string, 0, len(deps))
	for _, dep := range deps {
		// Split at the first occurrence of any version specifier
		splitIndex := len(dep)
		for i, char := range dep {
			if char == '=' || char == '<' || char == '>' || char == '!' || char == '~' {
				splitIndex = i
				break
			}
		}
		result = append(result, dep[:splitIndex])
	}
	return result
}

func addUpdatedDeps(directory string, mainDeps []string, devDeps []string, saveExact bool) error {
	if len(mainDeps) > 0 {
		log.Info().
			Int("numMainDeps", len(mainDeps)).
			Msg("Adding main dependencies back using 'uv add'")
		if err := addMainDeps(directory, mainDeps); err != nil {
			return fmt.Errorf("failed to add main dependencies: %w", err)
		}

		log.Info().Msg("Main dependencies added successfully")
	}

	if len(devDeps) > 0 {
		log.Info().
			Int("numDevDeps", len(devDeps)).
			Msg("Adding dev dependencies back using 'uv add --dev'")
		if err := addDevDeps(directory, devDeps); err != nil {
			return fmt.Errorf("failed to add dev dependencies: %w", err)
		}
		log.Info().Msg("Dev dependencies added successfully")
	}

	if saveExact {
		allDeps := make([]string, 0, len(mainDeps)+len(devDeps))
		allDeps = append(allDeps, mainDeps...)
		allDeps = append(allDeps, devDeps...)
		// Replace ">=" with "==" in dependencies in pyproject.toml "main" and "devDependencies"
		if err := makeVersionsExact(directory, allDeps); err != nil {
			return fmt.Errorf("failed to make versions exact: %w", err)
		}

		log.Info().Msg("Updated pyproject.toml to save exact versions of dependencies")
	}

	// run "uv sync" to generate lock file
	if err := updateLockFile(directory); err != nil {
		return fmt.Errorf("failed to sync uv environment: %w", err)
	}

	return nil
}

func makeVersionsExact(directory string, deps []string) error {
	pyprojectPath := path.Join(directory, "pyproject.toml")
	data, err := getPyProjectTomlData(pyprojectPath)
	if err != nil {
		return fmt.Errorf("failed to read pyproject.toml: %w", err)
	}

	tomlData := string(data)
	log.Info().
		Int("numDeps", len(deps)).
		Msg("Making version exact for dependency")
	for _, dep := range deps {
		tomlData = strings.ReplaceAll(tomlData, dep+">=", dep+"==")
	}

	if err := os.WriteFile(pyprojectPath, []byte(tomlData), 0644); err != nil {
		return fmt.Errorf("failed to write updated pyproject.toml: %w", err)
	}

	log.Info().Msg("Updated dependencies to exact versions in pyproject.toml")
	return nil

}

func addMainDeps(directory string, deps []string) error {
	cmd := []string{"uv", "add", "--directory", directory}
	cmd = append(cmd, deps...)
	return executeCommand(cmd)
}

func addDevDeps(directory string, deps []string) error {
	cmd := []string{"uv", "add", "--dev", "--directory", directory}
	cmd = append(cmd, deps...)
	return executeCommand(cmd)
}
