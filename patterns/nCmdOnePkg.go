package patterns

type NCmdOnePkg struct {
	NCmd
	OnePkg
}

func NewNCmdOnePkg(mode string, info []StructureInfo) NCmdOnePkg {
	nCmd := NewNCmd(mode, info)
	onePkg := NewOnePkg(mode, info)

	return NCmdOnePkg{nCmd, onePkg}
}

func (s NCmdOnePkg) Test() bool {
	return s.NCmd.Test() && s.OnePkg.Test()
}

func (s NCmdOnePkg) Spawn() (Fly, error) {
	nCmdSpawn, err := s.NCmd.Spawn()

	if err != nil {
		return Fly{}, err
	}
	onePkgSpawn, err := s.OnePkg.Spawn()

	if err != nil {
		return Fly{}, err
	}

	nCmdSpawn.Programs = append(nCmdSpawn.Programs, onePkgSpawn.Programs...)

	return nCmdSpawn, nil
}
