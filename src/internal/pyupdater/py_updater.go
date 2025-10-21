package pyupdater

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"path"
)

type _PyProjectToml struct {
	// Relevant fields are "project.dependencies" and "dependency-groups.dev"
	Project struct {
		Dependencies []string `toml:"dependencies"`
	}
	DependencyGroups map[string][]string `toml:"dependency-groups"`
}

func (p _PyProjectToml) GetMainDependencies() []string {
	return p.Project.Dependencies
}

func (p _PyProjectToml) GetDevDependencies() []string {
	return p.DependencyGroups["dev"]
}

func UpdatePackages(directory string, saveExact bool) error {
	// Get the path to the pyproject.toml file
	pyprojectPath := path.Join(directory, "pyproject.toml")
	data, err := getPyProjectTomlData(pyprojectPath)
	if err != nil {
		return fmt.Errorf("failed to get pyproject.toml data: %w", err)
	}

	pyproject, err := parsePyProjectToml(data)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", pyprojectPath, err)
	}

	log.Debug().
		Int("numDependencies", len(pyproject.GetMainDependencies())).
		Int("numDevDependencies", len(pyproject.GetDevDependencies())).
		// Str("firstMainDependency", pyproject.GetMainDependencies()[0]).
		Msg("Parsed dependencies from pyproject.toml")

	mainDeps := withoutVersion(pyproject.GetMainDependencies())
	devDeps := withoutVersion(pyproject.GetDevDependencies())
	log.Trace().
		Strs("mainDependencies", mainDeps).
		Strs("devDependencies", devDeps).
		Msg("Dependencies without version specifiers")

	// Now, remove lock file
	if err := removeLockFile(directory); err != nil {
		return err
	}

	// Then remove the main and dev dependencies from pyproject.toml
	if err := removeDependenciesFromTomlFile(data, pyprojectPath); err != nil {
		return err
	}

	// Now, add dependencies back using "uv add" and "uv add --dev"
	if err := addUpdatedDeps(directory, mainDeps, devDeps, saveExact); err != nil {
		return fmt.Errorf("failed to add updated dependencies: %w", err)
	}

	return nil
}
