package bitmart

import (
	"encoding/json"
	"fmt"
	slinkymath "github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	"math"
)

type (
	// Operation is the operation to perform. This is used to construct subscription messages
	// when initially connecting to the websocket. This can later be extended to support
	// other operations.
	Operation string
	// Channel is the channel to subscribe to. The channel is used to determine the type of
	// price data that we want. This can later be extended to support other channels. Currently,
	// only the index tickers (spot markets) channel is supported.
	Channel string
	// EventType is the event type. This is the expected event type that we want to receive
	// from the websocket. The event types pertain to subscription events.
	EventType string
)

const (
	// OperationSubscribe is the operation to subscribe to a channel.
	OperationSubscribe Operation = "subscribe"
)

const (
	// TickersChannel is the channel for tickers. This includes the spot price of the instrument.
	//
	TickersChannel Channel = "spot/ticker"
)

const (
	// EventSubscribe is the event denoting that we have successfully subscribed to a channel.
	EventSubscribe EventType = "subscribe"
	// EventTickers is the event for tickers. By default, this field will not be populated
	// in a properly formatted message. So we set the default value to an empty string.
	EventTickers EventType = ""
	// EventError is the event for an error.
	EventError EventType = "error"
)

// BaseMessage is utilized to determine the type of message that was received.
type BaseMessage struct {
	// Event is the event that occurred.
	Event string `json:"event" validate:"required"`
}

type SubscribeRequestMessage struct {
	// Operation is the operation to perform.
	Operation string `json:"op" validate:"required"`

	// Arguments is the list of arguments for the operation.
	Arguments []string `json:"args" validate:"required"`
}

// NewSubscribeToTickersRequestMessage returns a new SubscribeRequestMessage for subscribing
// to the tickers channel.
func (h *WebSocketHandler) NewSubscribeToTickersRequestMessage(
	instruments []string,
) ([]handlers.WebsocketEncodedMessage, error) {
	numInstruments := len(instruments)
	if numInstruments == 0 {
		return nil, fmt.Errorf("instruments cannot be empty")
	}
	numBatches := int(math.Ceil(float64(numInstruments) / float64(h.ws.MaxSubscriptionsPerBatch)))
	msgs := make([]handlers.WebsocketEncodedMessage, numBatches)
	for i := 0; i < numBatches; i++ {
		// Get the instruments for this batch.
		start := i * h.ws.MaxSubscriptionsPerBatch
		end := slinkymath.Min((i+1)*h.ws.MaxSubscriptionsPerBatch, numInstruments)
		batch := instruments[start:end]
		args := make([]string, 0)
		for _, instrument := range batch {
			args = append(args, fmt.Sprintf("%s:%s", TickersChannel, instrument))
		}

		bz, err := json.Marshal(
			SubscribeRequestMessage{
				Operation: string(OperationSubscribe),
				Arguments: args,
			},
		)
		if err != nil {
			return msgs, err
		}
		msgs[i] = bz
	}

	return msgs, nil
}

type SubscribeResponseMessage struct {
	// Event is the event that occurred.
	Event string `json:"event" validate:"required"`
	Topic string `json:"topic" validate:"required"`

	ErrorCode    string `json:"errorCode" `
	ErrorMessage string `json:"errorMessage"`
}

type TickersResponseMessage struct {
	Table string `json:"table" validate:"required"`
	// Data is the list of index ticker data.
	Data []IndexTicker `json:"data" validate:"required"`
}

// IndexTicker is the index ticker data.
type IndexTicker struct {
	// ID is the instrument ID.
	Symbol string `json:"symbol" validate:"required"`

	// LastPrice is the last price.
	LastPrice string `json:"last_price" validate:"required"`
}
