package detect

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// DetectStack is used to detect the stack of a project (recursively)
func DetectStack(projectRoot string) (*DetectedStack, error) {
	stack := &DetectedStack{}

	// Først: prøv roden, så root får lov at “vinde”
	if err := detectFromComposer(projectRoot, stack); err != nil {
		return nil, err
	}
	if err := detectFromPackageJson(projectRoot, stack); err != nil {
		return nil, err
	}
	if err := detectFromPackageLockJson(projectRoot, stack); err != nil {
		return nil, err
	}

	ignoredDirs := map[string]bool{
		"node_modules": true,
		"composer":     true,
		"vendor":       true,
	}

	err := filepath.WalkDir(projectRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// hvis der er en random permission-fejl et sted, så spring den bare over
			return nil
		}

		// spring root selv over – den er allerede kørt
		if path == projectRoot {
			return nil
		}

		if d.IsDir() {
			name := d.Name()

			// skip dot-foldere: .git, .idea, .vscode, ...
			if strings.HasPrefix(name, ".") {
				return fs.SkipDir
			}

			// skip bestemte foldere
			if ignoredDirs[name] {
				return fs.SkipDir
			}

			return nil
		}

		switch d.Name() {
		case "composer.json":
			_ = detectFromComposer(filepath.Dir(path), stack)
		case "package.json":
			_ = detectFromPackageJson(filepath.Dir(path), stack)
		case "package-lock.json":
			_ = detectFromPackageLockJson(filepath.Dir(path), stack)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return stack, nil
}
