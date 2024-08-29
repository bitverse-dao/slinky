package bitverse

import (
	"encoding/json"
	"fmt"

	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
)

type (
	// Operation is the operation to perform. This is used to construct subscription messages
	// when initially connecting to the websocket. This can later be extended to support
	// other operations.
	Operation string

	// Channel is the channel to subscribe to. The channel is used to determine the type of
	// price data that we want. This can later be extended to support other channels.
	Channel string
)

const (
	// OperationSubscribe is the operation to subscribe to a channel.
	OperationSubscribe Operation = "subscribe"

	OperationCommandResp Operation = "COMMAND_RESP"

	OperationSnapshot Operation = "snapshot"

	OperationPing Operation = "ping"

	// TickerChannel is the channel for spot price updates.
	TickerChannel Channel = "tickers"

	// MaxArgsPerRequest is the maximum amount of arguments that can be made for a single request to the ByBit WS API.
	MaxArgsPerRequest = 10
)

type BaseRequest struct {
	Op string `json:"op"`
}

// SubscriptionRequest is a request to the server to subscribe to ticker updates for currency pairs.
//
// Example:
//
//	{
//	   "req_id": "test", // optional
//	   "op": "subscribe",
//	   "args": [
//	       "orderbook.1.BTCUSDT",
//	       "publicTrade.BTCUSDT",
//	       "orderbook.1.ETHUSDT"
//	   ]
//	}
//	{
//		"op": "subscribe",
//		"id": "1704449555000",
//		"args": [
//			"tickers.BTC-USD",
//			"tickers.SATS-USD"
//		]
//	}

type SubscriptionRequest struct {
	BaseRequest
	ReqID string   `json:"id"`
	Args  []string `json:"args"`
}

// NewSubscriptionRequestMessage creates subscription messages corresponding to the provided tickers.
// If the number of tickers is greater than 10, the requests will be broken into 10-ticker messages.
func NewSubscriptionRequestMessage(tickers []string) ([]handlers.WebsocketEncodedMessage, error) {
	numTickers := len(tickers)
	if numTickers == 0 {
		return nil, fmt.Errorf("tickers cannot be empty")
	}

	messages := make([]handlers.WebsocketEncodedMessage, len(tickers))

	for i, ticker := range tickers {

		bz, err := json.Marshal(
			SubscriptionRequest{
				BaseRequest: BaseRequest{
					Op: string(OperationSubscribe),
				},
				Args: []string{ticker},
			},
		)
		if err != nil {
			return messages, fmt.Errorf("unable to marshal message: %w", err)
		}

		messages[i] = bz

	}

	return messages, nil
}

// HeartbeatPing is the ping sent to the server.
//
// Example:
//
//	{
//	   "req_id": "100010",
//	   "op": "ping"
//	}
type HeartbeatPing struct {
	BaseRequest
}

// NewHeartbeatPingMessage returns the encoded message for sending a heartbeat message to a peer.
func NewHeartbeatPingMessage() ([]handlers.WebsocketEncodedMessage, error) {
	bz, err := json.Marshal(
		HeartbeatPing{
			BaseRequest{
				Op: string(OperationPing),
			},
		},
	)

	return []handlers.WebsocketEncodedMessage{bz}, err
}

// HeartbeatPong is the pong sent back from the server after a ping.
//
// Example:
//
//	{
//	   "success": true,
//	   "ret_msg": "pong",
//	   "conn_id": "0970e817-426e-429a-a679-ff7f55e0b16a",
//	   "op": "ping"
//	}
type HeartbeatPong struct {
	Op   string   `json:"op"`
	Args []string `json:"args"`
}

// BaseMessage represents a base message. This is used to determine the type of message
// that was received.
type BaseMessage struct {
	// Type is the type of message.
	Type string `json:"type"`
}

// SubscriptionResponse is the response for a subscribe event.
type SubscriptionResponse struct {
	ReqID   string `json:"id"`
	Success bool   `json:"success"`
	ConnId  string `json:"conn_id"`
	Data    struct {
		FailTopics    []string `json:"fail_topics"`
		SuccessTopics []string `json:"success_topics"`
	} `json:"Data"`
	Type   string `json:"type"`
	RetMsg string `json:"ret_msg"`
}

// TickerUpdateMessage is the update sent for a subscribed ticker on the ByBit websocket API.
type TickerUpdateMessage struct {
	Id    string           `json:"id"`
	Ts    int64            `json:"ts"`
	Type  string           `json:"type"`
	Topic string           `json:"topic"`
	Data  TickerUpdateData `json:"data"`
}

// TickerUpdateData is the data stored inside a ticker update message.
type TickerUpdateData struct {
	Symbol          string `json:"symbol"`
	LastPrice       string `json:"lastPrice"`
	Price24hPercent string `json:"price24hPercent"`
	Price24hChange  string `json:"price24hChange"`
	OpenPrice       string `json:"openPrice"`
	HighPrice       string `json:"highPrice"`
	LowPrice        string `json:"lowPrice"`
	IndexPrice      string `json:"indexPrice"`
	OraclePrice     string `json:"oraclePrice"`
	Volume24h       string `json:"volume24h"`
	OpenTime        int64  `json:"openTime"`
	CloseTime       int64  `json:"closeTime"`
}
