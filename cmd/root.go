package cmd

import (
	"github.com/Duchadian/structurizr-export-cli/internal"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "structurizr-export-cli <structurizr url>",
	Short: "Export PNG images from a structurizr instance",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		rodUrl, _ := cmd.Flags().GetString("rod-remote")
		exportDir, _ := cmd.Flags().GetString("export-dir")
		structurizrUrl := args[0]

		internal.ExtractImages(structurizrUrl, rodUrl, exportDir)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("rod-remote", "", "remote rod instance to run the browser")
	rootCmd.Flags().String("export-dir", "export", "directory to export the images to")
}
