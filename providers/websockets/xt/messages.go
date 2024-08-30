package xt

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	slinkymath "github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
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
	TickersChannel Channel = "ticker"
)

const (
	// EventSubscribe is the event denoting that we have successfully subscribed to a channel.
	EventSubscribe EventType = "subscribe"
	// EventTickers is the event for tickers. By default, this field will not be populated
)

// BaseMessage is utilized to determine the type of message that was received.
type BaseMessage struct {
	// Event is the event that occurred.
	Event string `json:"event" validate:"required"`
}

type SubscribeRequestMessage struct {
	ID     string   `json:"id" validate:"required"`
	Method string   `json:"method" validate:"required"`
	Params []string `json:"params" validate:"required"`
}

// NewSubscribeRequest returns a new SubscribeRequest encoded message for the given symbols.
func (h *WebSocketHandler) NewSubscribeRequest(symbols []string) ([]handlers.WebsocketEncodedMessage, error) {
	numSymbols := len(symbols)
	if numSymbols == 0 {
		return nil, fmt.Errorf("cannot attach payload of 0 length")
	}

	numBatches := int(math.Ceil(float64(numSymbols) / float64(h.ws.MaxSubscriptionsPerBatch)))
	msgs := make([]handlers.WebsocketEncodedMessage, numBatches)
	for i := 0; i < numBatches; i++ {
		// Get the symbols for the batch.
		start := i * h.ws.MaxSubscriptionsPerBatch
		end := slinkymath.Min((i+1)*h.ws.MaxSubscriptionsPerBatch, numSymbols)
		batch := symbols[start:end]
		params := make([]string, 0)
		for _, b := range batch {
			// ticker@btc_usdt
			params = append(params, fmt.Sprintf("%s@%s", string(TickersChannel), strings.ToLower(b)))
		}
		bz, err := json.Marshal(SubscribeRequestMessage{
			Method: string(OperationSubscribe),
			ID:     strconv.Itoa(time.Now().UTC().Second()),
			Params: params,
		})
		if err != nil {
			return msgs, err
		}
		msgs[i] = bz
	}
	return msgs, nil
}

type SubscribeResponseMessage struct {
	ID int `json:"id,omitempty"`

	// Code is the error code.
	Code int `json:"code,omitempty"`

	// Message is the error message. Note that the field will be populated with the same exact
	// initial message that was sent to the websocket.
	Message string `json:"msg,omitempty"`
}

type TickersResponseMessage struct {
	Topic string `json:"topic"`
	// Data is the list of index ticker data.
	Data []IndexTicker `json:"data" validate:"required"`
}

// IndexTicker is the index ticker data.
type IndexTicker struct {
	// ID is the instrument ID.
	Symbol string `json:"s" validate:"required"`

	// LastPrice is the last price.
	LastPrice string `json:"c" validate:"required"`
}
