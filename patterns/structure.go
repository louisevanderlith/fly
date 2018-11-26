package patterns

import (
	"os"
	"path/filepath"
)

func generateConfig(path, mode string) Fly {
	result := Fly{
		Env: environment{
			Bin:  "./bin",
			Mode: mode,
		},
	}

	apps, err := scanForMain(path)

	if err != nil {
		panic(err)
	}

	for name, appPath := range apps {
		p := Program{
			Name:     name,
			Path:     appPath,
			Priority: 0,
			Play:     true,
		}

		result.Programs = append(result.Programs, p)
	}

	return result
}

func scanForMain(basePath string) (map[string]string, error) {
	result := make(map[string]string)

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if info.Name() == "main.go" {
			dir, _ := filepath.Split(path)
			name := filepath.Base(dir)
			result[name] = dir
		}

		return nil
	})

	return result, err
}
