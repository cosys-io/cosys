package common

import (
	"github.com/spf13/cobra"
	"log"
)

type Command func(*Cosys) *cobra.Command

func (c Command) String() string {
	return c(nil).Name()
}

func rootCmd(cosys *Cosys) *cobra.Command {
	root := &cobra.Command{}

	root.AddCommand(serveCmd(cosys))

	return root
}

func serveCmd(cosys *Cosys) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cosys.Bootstrap(); err != nil {
				log.Fatal(err)
			}

			if err := cosys.Server().Start(); err != nil {
				log.Fatal(err)
			}

			if err := cosys.Destroy(); err != nil {
				log.Fatal(err)
			}
		},
	}
}
