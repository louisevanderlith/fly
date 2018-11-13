package patterns

//OneCmd is [1]-Command
type OneCmd struct {
	mode string
	info []StructureInfo
}

func NewOneCmd(mode string, info []StructureInfo) OneCmd {
	return OneCmd{mode, info}
}

//Test should pass when the is only one main.go in the folder.
func (s OneCmd) Test() bool {
	mainsFound := 0

	for _, v := range s.info {
		if v.HasMain {
			mainsFound++
		}
	}

	return mainsFound == 1
}

func (s OneCmd) Spawn() (Fly, error) {
	result := Fly{
		Env: environment{
			Bin:  "$GOPATH/bin",
			Mode: s.mode,
		},
	}

	for _, v := range s.info {
		prog := Program{
			Name:     v.Name,
			Play:     true,
			Priority: 1,
			Type:     Cmd,
		}

		result.Programs = append(result.Programs, prog)
	}

	return result, nil
}
