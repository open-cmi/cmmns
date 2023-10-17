package subcommands

type HelpCommand struct {
}

func (c *HelpCommand) Synopsis() string {
	return "print usage"
}

func (c *HelpCommand) Run() error {
	Usage()
	return nil
}

func init() {
	RegisterCommand("help", &HelpCommand{})
}
