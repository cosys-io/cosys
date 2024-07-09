package cms

import (
	"github.com/cosys-io/cosys/common"
	"log"
)

func init() {
	if err := common.RegisterCommand(rootCmd); err != nil {
		log.Fatal(err)
	}
}
