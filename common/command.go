package common

import (
	"github.com/spf13/cobra"
	"log"
)

// Command takes in a cosys instance and returns a command.
type Command func(*Cosys) *cobra.Command

// String returns the name of the command, i.e. the word used to run the command in the cli.
func (c Command) String() string {
	return c(nil).Name()
}

// serveCmd is the command to start the server in production mode.
func serveCmd(cosys *Cosys) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the server in production mode",
		Run: func(cmd *cobra.Command, args []string) {
			cosys.SetEnvironment(Prod)

			for err := range cosys.startServer() {
				log.Print(err)
			}
		},
	}
}

// devCmd is the command to start the server in development mode.
func devCmd(cosys *Cosys) *cobra.Command {
	return &cobra.Command{
		Use:   "dev",
		Short: "Start the server in development mode",
		Run: func(cmd *cobra.Command, args []string) {
			cosys.SetEnvironment(Dev)

			for err := range cosys.startServer() {
				log.Print(err)
			}
		},
	}
}

// testCmd is the command to start the server in test mode.
func testCmd(cosys *Cosys) *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Start the server in test mode",
		Run: func(cmd *cobra.Command, args []string) {
			cosys.SetEnvironment(Test)

			for err := range cosys.startServer() {
				log.Print(err)
			}
		},
	}
}
