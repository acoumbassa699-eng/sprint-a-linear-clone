package cli

import (
	"github.com/optimus-ide-collab/serpent"
)

func (r *RootCmd) tasksCommand() *serpent.Command {
	cmd := &serpent.Command{
		Use:     "task",
		Aliases: []string{"tasks"},
		Short:   "Manage tasks",
		Handler: func(i *serpent.Invocation) error {
			return i.Command.HelpHandler(i)
		},
		Children: []*serpent.Command{
			r.taskCreate(),
			r.taskDelete(),
			r.taskList(),
			r.taskLogs(),
			r.taskPause(),
			r.taskResume(),
			r.taskSend(),
			r.taskStatus(),
		},
	}
	return cmd
}
