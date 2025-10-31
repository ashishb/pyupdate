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

func removeDependenciesFromTomlFile(data []byte, deps []string, pyprojectPath string) error {
	tomlData := string(data)
	// We cannot use toml.Unmarshal and Marshal here because it changes the formatting of the file.
	numRemoved := 0
	for _, dep := range deps {
		if strings.Count(tomlData, `"`+dep+`"`) > 1 && strings.Count(tomlData, `'`+dep+`'`) > 1 {
			return fmt.Errorf("dependency %s appears multiple times in pyproject.toml", dep)
		}

		if strings.Count(tomlData, dep) == 0 {
			return fmt.Errorf("dependency %s not found in pyproject.toml", dep)
		}

		val1 := fmt.Sprintf(`"%s",`, dep)
		val2 := fmt.Sprintf(`'%s',`, dep)
		val3 := fmt.Sprintf(`"%s"`, dep)
		for _, val := range []string{val1, val2, val3} {
			if strings.Contains(tomlData, val) {
				tomlData = strings.Replace(tomlData, val, "", 1)
				numRemoved++
				break
			}
		}
	}

	if err := os.WriteFile(pyprojectPath, []byte(tomlData), 0644); err != nil {
		return fmt.Errorf("failed to write updated pyproject.toml: %w", err)
	}

	log.Info().
		Int("numRemovedDeps", numRemoved).
		Msg("Removed main + dev dependencies from pyproject.toml")

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

func addUpdatedDeps(directory string, mainDeps []string, devDeps []string, optionalDeps map[string][]string, saveExact bool) error {
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

	if len(optionalDeps) > 0 {
		log.Info().
			Int("numOptionalDeps", len(optionalDeps)).
			Msg("Adding optional dependencies back using 'uv add --optional'")
		if err := addOptionalDeps(directory, optionalDeps); err != nil {
			return fmt.Errorf("failed to add optional dependencies: %w", err)
		}
		log.Info().Msg("Optional dependencies added successfully")
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

func addOptionalDeps(directory string, deps map[string][]string) error {
	for group, groupDeps := range deps {
		cmd := []string{"uv", "add", "--optional", group, "--directory", directory}
		cmd = append(cmd, groupDeps...)
		if err := executeCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}
