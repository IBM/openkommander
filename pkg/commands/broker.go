package commands

import (
	"github.com/IBM/openkommander/pkg/broker"
)

func brokerInfoCommand() {
	broker.GetInfo()
}
