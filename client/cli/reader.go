package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/afrozalm/minimess/client/connHandler"
	"github.com/afrozalm/minimess/constants"
	"github.com/afrozalm/minimess/message"
)

func Run(h *connHandler.Handler) {
	reader := bufio.NewReader(os.Stdin)
	messageHolder := new(message.Message)
	messageHolder.Uid = h.Uid

	for {
		fmt.Print("> ")
		inputString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		inputString = strings.Trim(inputString, "\n ")
		split := strings.SplitN(inputString, " ", 2)
		cmd, rest := split[0], split[1]
		switch cmd {
		case constants.SUBSCRIBE:
			messageHolder.Type = constants.SUBSCRIBE
			messageHolder.Topic = rest
			h.Send <- messageHolder
		case constants.UNSUBSCRIBE:
			messageHolder.Type = constants.UNSUBSCRIBE
			messageHolder.Topic = rest
			h.Send <- messageHolder
		case constants.CHAT:
			split = strings.SplitN(rest, " ", 2)
			topic, text := split[0], split[1]
			messageHolder.Type = constants.CHAT
			messageHolder.Topic = topic
			messageHolder.Text = text
			h.Send <- messageHolder
		default:
			fmt.Fprintln(os.Stderr, "bad message type")
		}
	}
}
