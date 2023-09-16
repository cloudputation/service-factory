package cli

import (
	"fmt"
)

func HelperMessage()() {
	helpMessage := `Usage: servicectl <command> [args]

Help:           Description:                    Usage:
    help        Get help menu                   servicectl help



Commands:       Description:                    Usage:
    config      Parse config and print          servicectl config [path]
    apply       Deploy service                  servicectl apply [path]
    agent       Run Service Factory server      servicectl agent
`

	fmt.Println(helpMessage)
}
