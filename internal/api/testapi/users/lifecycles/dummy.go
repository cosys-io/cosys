package lifecycles

import (
	"log"

	"github.com/cosys-io/cosys/internal/common"
)

var Dummy = common.AfterFindOne(func(event *common.Event) error {
	log.Print(event.Result)
	return nil
})
