package coprocessor

import (
	"encoding/json"
)

type FinishRequest struct {
	Status string `json:"status"`
}

type FinishResponse struct {
	Type string          `json:"request_type"`
	Data json.RawMessage `json:"data"`
}

type AdvanceResponse struct {
	Metadata struct {
		ChainId     uint64 `json:"chain_id"`
		TaskManager string `json:"task_manager"`
		MsgSender   string `json:"msg_sender"`
		BlockHash   string `json:"block_hash"`
		BlockNumber uint64 `json:"block_number"`
		Timestamp   uint64 `json:"timestamp"`
		PrevRandao  string `json:"prev_randao"`
	} `json:"metadata"`
	Payload string `json:"payload"`
}

type NoticeRequest struct {
	Payload string `json:"payload"`
}

type ExceptionRequest struct {
	Payload string `json:"payload"`
}

type IndexResponse struct {
	Index uint64 `json:"index"`
}
