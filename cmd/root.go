package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "trader",
	Short: "Investor Support Tool",
	Long: `Investor Support Tool

  Designed to assist investors of all levels in making 
  informed decisions and optimizing their portfolios. 
  With an intuitive interface and access to real-time data, 
  you will have the necessary information to track the market and 
  analyze opportunities with confidence.
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
