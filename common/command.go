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
		Short: "Start the server in production mode",
		Run: func(cmd *cobra.Command, args []string) {
			cosys.SetEnvironment(Prod)

			if err := cosys.Bootstrap(); err != nil {
				log.Fatal(err)
			}

			go func() {
				if err := cosys.Server().Start(); err != nil {
					log.Fatal(err)
				}
			}()

			<-cosys.ShutdownChannel()

			if err := cosys.Destroy(); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func devCmd(cosys *Cosys) *cobra.Command {
	return &cobra.Command{
		Use:   "dev",
		Short: "Start the server in development mode",
		Run: func(cmd *cobra.Command, args []string) {
			cosys.SetEnvironment(Dev)

			if err := cosys.Bootstrap(); err != nil {
				log.Fatal(err)
			}

			go func() {
				if err := cosys.Server().Start(); err != nil {
					log.Fatal(err)
				}
			}()

			<-cosys.ShutdownChannel()

			if err := cosys.Destroy(); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func testCmd(cosys *Cosys) *cobra.Command {
	return &cobra.Command{
		Use:   "dev",
		Short: "Start the server in test mode",
		Run: func(cmd *cobra.Command, args []string) {
			cosys.SetEnvironment(Test)

			if err := cosys.Bootstrap(); err != nil {
				log.Fatal(err)
			}

			go func() {
				if err := cosys.Server().Start(); err != nil {
					log.Fatal(err)
				}
			}()

			<-cosys.ShutdownChannel()

			if err := cosys.Destroy(); err != nil {
				log.Fatal(err)
			}
		},
	}
}
