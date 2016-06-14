package commands

type ASGCreatorCommand struct {
	Create CreateCommand `command:"create" description:"Create default ASGs"`
}

var ASGCreator ASGCreatorCommand
