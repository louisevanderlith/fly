package patterns

type OneCmdNPkg struct {
	OneCmd
	NPkg
}

func NewOneCmdNPkg(mode string, info []StructureInfo) OneCmdNPkg {
	oneCmd := NewOneCmd(mode, info)
	nPkg := NewNPkg(mode, info)

	return OneCmdNPkg{oneCmd, nPkg}
}

func (s OneCmdNPkg) Test() bool {
	return s.OneCmd.Test() && s.NPkg.Test()
}

func (s OneCmdNPkg) Spawn() (Fly, error) {
	oneCmdSpawn, err := s.OneCmd.Spawn()

	if err != nil {
		return Fly{}, err
	}

	nPkgSpawn, err := s.NPkg.Spawn()

	if err != nil {
		return Fly{}, err
	}

	oneCmdSpawn.Programs = append(oneCmdSpawn.Programs, nPkgSpawn.Programs...)

	return oneCmdSpawn, nil
}
