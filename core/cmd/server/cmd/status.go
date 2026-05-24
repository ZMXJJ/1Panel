package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/1Panel-dev/1Panel/core/utils/controller"
	"github.com/1Panel-dev/1Panel/core/utils/platform"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show panel service status",
	RunE: func(cmd *cobra.Command, args []string) error {
		return printPanelStatus(cmd.OutOrStdout())
	},
}

func printPanelStatus(out io.Writer) error {
	if platform.IsDarwin() {
		printLaunchdStatus(out, "core", "dev.oasis.core")
		printLaunchdStatus(out, "agent", "dev.oasis.agent")
		return nil
	}
	printLinuxServiceStatus(out, "core", "1panel-core")
	printLinuxServiceStatus(out, "agent", "1panel-agent")
	return nil
}

func printLaunchdStatus(out io.Writer, name, label string) {
	state := "inactive"
	if err := exec.Command("launchctl", "print", fmt.Sprintf("gui/%d/%s", os.Getuid(), label)).Run(); err == nil {
		state = "active"
	}
	fmt.Fprintf(out, "%s\t%s\t%s\n", name, state, label)
}

func printLinuxServiceStatus(out io.Writer, name, service string) {
	state := "inactive"
	if active, err := controller.CheckActive(service); err == nil && active {
		state = "active"
	}
	fmt.Fprintf(out, "%s\t%s\t%s\n", name, state, service)
}
