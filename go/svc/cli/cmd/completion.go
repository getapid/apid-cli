package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var completionFilepath string

var completionCmd = &cobra.Command{
	Use:       "completion",
	Short:     "Generates bash completion scripts",
	ValidArgs: []string{"bash", "zsh", "powershell"},
	Args:      cobra.ExactArgs(1),
	Long: `To load completion run

. <(apid completion <shell>)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(apid completion <shell>)
`,
	Run: func(cmd *cobra.Command, args []string) {
		destination := os.Stdout

		if completionFilepath != "" {
			f, err := os.Open(completionFilepath)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			destination = f
		}

		var err error
		switch args[1] {
		case "bash":
			err = RootCmd.GenBashCompletion(destination)
		case "zsh":
			err = RootCmd.GenZshCompletion(destination)
		case "powershell":
			err = RootCmd.GenPowerShellCompletion(destination)
		}
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(completionCmd)
	completionCmd.Flags().StringVarP(&completionFilepath, "file", "f", "", "the file to output the completion script to (defaults to stdout)")
}