package patterns

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//Structurer will pass (info []StructureInfo) to a NewPattern, for testing
type Structurer interface {
	Test() bool
	Spawn() (Fly, error)
}

func registeredStructures(mode string, info []StructureInfo) []Structurer {
	var structures []Structurer

	structures = append(structures, NewNCmdNPkg(mode, info))
	structures = append(structures, NewNCmdOnePkg(mode, info))
	structures = append(structures, NewOneCmdNPkg(mode, info))
	structures = append(structures, NewOneCmdOnePkg(mode, info))
	structures = append(structures, NewNCmd(mode, info))
	structures = append(structures, NewOneCmd(mode, info))
	structures = append(structures, NewNPkg(mode, info))
	structures = append(structures, NewOnePkg(mode, info))

	return structures
}

func generateConfig(path, mode string) (Fly, error) {
	info, err := ScanFolder(path)

	if err != nil {
		return Fly{}, err
	}

	for _, v := range registeredStructures(mode, info) {
		if v.Test() {
			return v.Spawn()
		}
	}

	return Fly{}, errors.New("no structure tests passed")
}

func ScanFolder(basePath string) ([]StructureInfo, error) {
	progMap := make(map[string][]string)

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		withoutName := strings.Replace(path, info.Name(), "", 1)

		progMap[withoutName] = append(progMap[withoutName], info.Name())

		return nil
	})

	if err != nil {
		return nil, err
	}

	return buildStructure(progMap), nil
}

func buildStructure(progMap map[string][]string) []StructureInfo {
	var result []StructureInfo

	for k, v := range progMap {
		prog := StructureInfo{
			HasGoFiles: false,
			HasMain:    false,
			Name:       lastFolder(k),
			Path:       k,
		}

		for i := 0; i < len(v); i++ {
			curr := v[i]

			if strings.Contains(curr, ".go") {
				prog.HasGoFiles = true
			}

			if curr == "main.go" {
				prog.HasMain = true
				break
			}
		}

		if prog.HasMain {
			result = append(result, prog)
		}
	}

	return result
}

func lastFolder(path string) string {
	parts := strings.Split(path, "/")

	//windows...
	if len(parts) == 1 {
		parts = strings.Split(path, "\\")
	}

	for i := (len(parts) - 1); i >= 0; i-- {
		curr := parts[i]

		if curr != "" {
			return curr
		}
	}

	return path
}

/*
   (x) 1 package
   (x) n packages
   (x) 1 command
   (x) 1 command and 1 package
   (x) n commands and 1 package
   (x) 1 command and n packages
   (x) n commands and n packages
*/
