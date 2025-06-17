package cmd

import (
	"fmt"
	"trader/internal/config"

	"github.com/spf13/cobra"
)

type RootCommand interface {
	GetCobraCommand() *cobra.Command
	Execute() error
}

type rootCommand struct {
	config config.Config
	// Commands
	rootCmd *cobra.Command
}

func (rc *rootCommand) getVersionCmd(cmd *cobra.Command, args []string) {
	fmt.Printf("%s\n", rc.config.GetVersion())
}

func (rc *rootCommand) setup() {
	getVersionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  `Show version`,
		Run:   rc.getVersionCmd,
	}
	rc.rootCmd.AddCommand(getVersionCmd)
}

func (rc *rootCommand) GetCobraCommand() *cobra.Command {
	return rc.rootCmd
}

func (rc *rootCommand) Execute() error {
	rc.setup()
	return rc.GetCobraCommand().Execute()
}

func NewRootCommand(config config.Config) RootCommand {
	return &rootCommand{
		config: config,
		rootCmd: &cobra.Command{
			Use:   "trader",
			Short: "Investor Support Tool",
			Long: `Investor Support Tool

  Designed to assist investors of all levels in making 
  informed decisions and optimizing their portfolios. 
  With an intuitive interface and access to real-time data, 
  you will have the necessary information to track the market and 
  analyze opportunities with confidence.
`,
		},
	}
}
