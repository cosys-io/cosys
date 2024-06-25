package common

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/logger/internal"
	"log"
)

type Logger struct {
}

func init() {
	if err := common.RegisterLogger("default", internal.Logger{}); err != nil {
		log.Fatal(err)
	}
}
