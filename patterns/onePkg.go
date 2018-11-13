package patterns

//OnePkg is [1]-Package
type OnePkg struct {
	mode string
	info []StructureInfo
}

func NewOnePkg(mode string, info []StructureInfo) OnePkg {
	return OnePkg{mode, info}
}

//Test should pass, only if the folder contains One package
func (s OnePkg) Test() bool {
	pkgCount := 0

	for _, v := range s.info {
		if !v.HasMain && v.HasGoFiles {
			pkgCount++
		}
	}

	return pkgCount == 1
}

func (s OnePkg) Spawn() (Fly, error) {
	result := Fly{
		Env: environment{
			Bin:  "$GOPATH/pkg",
			Mode: s.mode,
		},
	}

	for _, v := range s.info {
		prog := Program{
			Name:     v.Name,
			Path:     v.Path,
			Play:     false,
			Priority: 0,
			Type:     Pkg,
		}

		result.Programs = append(result.Programs, prog)
	}

	return result, nil
}
