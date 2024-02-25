package cmd

import (
	"fmt"
	"github.com/cmd-tools/gtfocli/constants"
	"github.com/cmd-tools/gtfocli/logger"
	"github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   constants.Update,
	Short: "Update gtfocli database.",
	Long:  "Download sources use but gtfocli to run offline search.",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init(constants.Update, IsDebug())
		err := downloadFiles()
		if err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

var yamlURLs = map[string]string{
	constants.GTFOBinsOutputDir: "github.com/GTFOBins/GTFOBins.github.io/_gtfobins",
	constants.LOLBASOutputDir:   "github.com/LOLBAS-Project/LOLBAS/yml/OSBinaries",
}

func downloadFiles() error {
	fmt.Println("Downloading...")

	for outputFolder, url := range yamlURLs {
		client := &getter.Client{
			Src:  url,
			Dst:  outputFolder,
			Pwd:  constants.MainOutputDir,
			Mode: getter.ClientModeDir,
		}

		logger.Logger.Infof("Downloading: %s into %s.", url, outputFolder)

		err := client.Get()

		if err != nil {
			fmt.Printf("failed to download folder %s: %s.\n", url, err.Error())
		}

		logger.Logger.Info("Done.")
	}

	fmt.Printf("Database updated!")

	return nil
}
