package daemon

import (
	"context"
	"time"

	"github.com/abiosoft/colima/cmd/root"
	"github.com/abiosoft/colima/config"
	"github.com/abiosoft/colima/daemon/process"
	"github.com/abiosoft/colima/daemon/process/gvproxy"
	"github.com/abiosoft/colima/daemon/process/vmnet"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:    "daemon",
	Short:  "daemon",
	Long:   `runner for background daemons.`,
	Hidden: true,
}

var startCmd = &cobra.Command{
	Use:   "start [profile]",
	Short: "start daemon",
	Long:  `start the daemon`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		config.SetProfile(args[0])
		ctx := cmd.Context()

		var processes []process.Process
		if daemonArgs.vmnet {
			processes = append(processes, vmnet.New())
		}
		if daemonArgs.gvproxy {
			processes = append(processes, gvproxy.New())
		}

		return start(ctx, processes)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop [profile]",
	Short: "stop daemon",
	Long:  `stop the daemon`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		config.SetProfile(args[0])

		// wait for 60 seconds
		timeout := time.Second * 60
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		return stop(ctx)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "status of the daemon",
	Long:  `status of the daemon`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		config.SetProfile(args[0])

		return status()
	},
}

var daemonArgs struct {
	vmnet    bool
	gvproxy  bool
	fsnotify bool

	verbose bool
}

func init() {
	root.Cmd().AddCommand(daemonCmd)

	daemonCmd.AddCommand(startCmd)
	daemonCmd.AddCommand(stopCmd)
	daemonCmd.AddCommand(statusCmd)

	startCmd.Flags().BoolVar(&daemonArgs.vmnet, "vmnet", false, "start vmnet")
	startCmd.Flags().BoolVar(&daemonArgs.gvproxy, "gvproxy", false, "start gvproxy")
	startCmd.Flags().BoolVar(&daemonArgs.fsnotify, "fsnotify", false, "start fsnotify")
}
