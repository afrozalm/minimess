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
	messageHolder := &message.Chat{
		UserID: h.Uid,
	}
	mode := constants.NORMAL

	for {
		fmt.Print("> ")
		inputString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		inputString = strings.Trim(inputString, "\n ")
		split := strings.Split(inputString, " ")
		switch {
		case len(split) == 0:
		case mode == constants.NORMAL && split[0] == constants.SUBSCRIBE:
			messageHolder.Type = constants.SUBSCRIBE
			rest := split[1:]
			if len(rest) == 0 {
				logBadMessageType(mode, "nothing to subscribe to")
			}
			for _, topicName := range rest {
				go func(topicName string) {
					sendSubscribe(h, topicName)
				}(topicName)
			}
		case mode == constants.NORMAL && split[0] == constants.UNSUBSCRIBE:
			messageHolder.Type = constants.UNSUBSCRIBE
			rest := split[1:]
			if len(rest) == 0 {
				logBadMessageType(mode, "")
			}
			for _, topicName := range rest {
				go func(topicName string) {
					sendUnsubscribe(h, topicName)
				}(topicName)
			}
		case mode == constants.NORMAL && split[0] == constants.ATTACH:
			topic := split[1]
			messageHolder.Type = constants.CHAT
			messageHolder.Topic = topic
			mode = constants.CHAT
		case mode == constants.CHAT:
			if split[0] == constants.DETACH {
				mode = constants.NORMAL
			} else {
				messageHolder.Text = inputString
				h.Send <- messageHolder
			}
		case mode == constants.NORMAL && split[0] == constants.QUIT:
			h.Done <- struct{}{}
			return
		default:
			logBadMessageType(mode, "")
		}
	}
}

func sendSubscribe(h *connHandler.Handler, topicName string) {
	h.Send <- &message.Chat{
		Type:   constants.SUBSCRIBE,
		UserID: h.Uid,
		Topic:  topicName,
	}
}

func sendUnsubscribe(h *connHandler.Handler, topicName string) {
	h.Send <- &message.Chat{
		Type:   constants.UNSUBSCRIBE,
		UserID: h.Uid,
		Topic:  topicName,
	}
}

func logBadMessageType(mode, msg string) {
	fmt.Fprintf(os.Stderr, "[mode:%s] bad message type - %s", mode, msg)
}
