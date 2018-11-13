package patterns

type NCmdNPkg struct {
	NCmd
	NPkg
}

func NewNCmdNPkg(mode string, info []StructureInfo) NCmdNPkg {
	nCmd := NewNCmd(mode, info)
	nPkg := NewNPkg(mode, info)

	return NCmdNPkg{nCmd, nPkg}
}

func (s NCmdNPkg) Test() bool {
	return s.NCmd.Test() && s.NPkg.Test()
}

func (s NCmdNPkg) Spawn() (Fly, error) {
	nCmdSpawn, err := s.NCmd.Spawn()

	if err != nil {
		return Fly{}, err
	}
	nPkgSpawn, err := s.NPkg.Spawn()

	if err != nil {
		return Fly{}, err
	}

	nCmdSpawn.Programs = append(nCmdSpawn.Programs, nPkgSpawn.Programs...)

	return nCmdSpawn, nil
}
