package patterns

//NCmd is [n]Commands
type NCmd struct {
	mode string
	info []StructureInfo
}

func NewNCmd(mode string, info []StructureInfo) NCmd {
	return NCmd{mode, info}
}

//Test should pass when there are more than one main.go in the folder
func (s NCmd) Test() bool {
	mainsFound := 0

	for _, v := range s.info {
		if v.HasMain {
			mainsFound++
		}
	}

	return mainsFound > 1
}

//Spawn
func (s NCmd) Spawn() (Fly, error) {
	result := Fly{
		Env: environment{
			Bin:  "./bin",
			Mode: s.mode,
		},
	}

	for _, v := range s.info {
		prog := Program{
			Name: v.Name,
			Path: v.Path,
			Play: true,
			Type: Cmd,
		}

		result.Programs = append(result.Programs, prog)
	}

	return result, nil
}
