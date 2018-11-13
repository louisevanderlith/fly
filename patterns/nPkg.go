package patterns

//NPkg is [n]Packages
type NPkg struct {
	mode string
	info []StructureInfo
}

func NewNPkg(mode string, info []StructureInfo) NPkg {
	return NPkg{mode, info}
}

//Test should pass when there are more than one package in the folder
func (s NPkg) Test() bool {
	pkgCount := 0

	for _, v := range s.info {
		if !v.HasMain && v.HasGoFiles {
			pkgCount++
		}
	}

	return pkgCount > 1
}

func (s NPkg) Spawn() (Fly, error) {
	result := Fly{
		Env: environment{
			Bin:  "$GOPATH/bin",
			Mode: s.mode,
		},
	}

	for _, v := range s.info {
		prog := Program{
			Name:     v.Name,
			Path:     v.Path,
			Priority: 0,
			Play:     false,
			Type:     Pkg,
		}

		result.Programs = append(result.Programs, prog)
	}

	return result, nil
}
