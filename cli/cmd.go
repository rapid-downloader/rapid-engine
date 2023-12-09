package main

import (
	"context"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/rapid-downloader/rapid/client"
	"github.com/spf13/cobra"
)

var cmds = make([]commandFunc, 0)

type commandFunc func(ctx context.Context, rapid client.Rapid) *cobra.Command

func registerCommand(cmd commandFunc) {
	cmds = append(cmds, cmd)
}

func executeCommands(ctx context.Context, rapid client.Rapid) {
	rootCmd := &cobra.Command{
		Use:   "rapid",
		Short: "Fetch and download",
		Long:  "Fetch and download a file from given url",
	}

	for _, command := range cmds {
		cmd := command(ctx, rapid)
		rootCmd.AddCommand(cmd)
	}

	rootCmd.Execute()
}

func download(ctx context.Context, rapid client.Rapid) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "download",
		Aliases: []string{"d"},
		Example: "rapid download <url> | rapid d <url>",
		Short:   "Download a file from the given url",
		Run: func(cmd *cobra.Command, args []string) {
			provider, _ := cmd.Flags().GetString("provider")
			if provider == "" {
				provider = "default"
			}

			s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
			s.Prefix = "Fetching url"
			s.Suffix = "\n"

			s.Start()
			defer s.Stop()

			url := args[0]
			request := client.Request{
				Url:      url,
				Provider: provider,
			}

			result, err := rapid.Fetch(request)
			if err != nil {
				log.Fatal(err)
				return
			}

			if err := rapid.Download(result.ID); err != nil {
				log.Fatal(err)
				return
			}

			store(result.ID, *result)
		},
	}

	return cmd
}

func init() {
	registerCommand(download)
}
