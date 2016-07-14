package cmd

import (
	"log"
	"os"
	"runtime"

	"github.com/zyfdegh/mdviewer/server"

	"github.com/spf13/cobra"
)

// RootCmd is root command of mdviewer
var RootCmd = &cobra.Command{
	Use:   "mdv",
	Short: "Mdviewer(mdv) is a markdown viewer, it displays markdown files in your browser.",
	Long:  "Mdviewer(mdv) is a markdown viewer, it opens your system default web browser and displays markdown files.",
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS != "linux" {
			log.Fatalln("Only linux is supported now.")
			return
		}

		// 0 flag, 1 arg
		// like "mdv readme.md"
		if cmd.Flags().NArg() == 1 {
			server.Serve(cmd.Flags().Args()[0])
		}

		// others
		cmd.SetArgs([]string{"--help"})
		if err := cmd.Execute(); err != nil {
			log.Fatalf("command arguments error: %v", err)
			os.Exit(-1)
		}
		return
	},
}
