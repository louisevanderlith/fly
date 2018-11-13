package patterns

type OneCmdOnePkg struct {
	OneCmd
	OnePkg
}

func NewOneCmdOnePkg(mode string, info []StructureInfo) OneCmdOnePkg {
	oneCmd := NewOneCmd(mode, info)
	onePkg := NewOnePkg(mode, info)

	return OneCmdOnePkg{oneCmd, onePkg}
}

func (s OneCmdOnePkg) Test() bool {
	return s.OneCmd.Test() && s.OnePkg.Test()
}

func (s OneCmdOnePkg) Spawn() (Fly, error) {
	oneCmdSpawn, err := s.OneCmd.Spawn()

	if err != nil {
		return Fly{}, err
	}

	onePkgSpawn, err := s.OnePkg.Spawn()

	if err != nil {
		return Fly{}, err
	}

	oneCmdSpawn.Programs = append(oneCmdSpawn.Programs, onePkgSpawn.Programs...)

	return oneCmdSpawn, nil
}
