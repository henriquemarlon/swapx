package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/henriquemarlon/swapx/configs"
	"github.com/henriquemarlon/swapx/pkg/coprocessor"
)

var (
	infolog = log.New(os.Stderr, "[ info ]  ", log.Lshortfile)
	errlog  = log.New(os.Stderr, "[ error ] ", log.Lshortfile)
)

func Handler(response *coprocessor.AdvanceResponse) error {
	_, err := configs.SetupInMemoryDB()
	if err != nil {
		errlog.Panicln("Failed to setup database", "error", err)
	}
	infolog.Println("Database setup successful")
	infolog.Println("Processing payload:", response)
	return nil
}

func main() {
	finish := coprocessor.FinishRequest{Status: "accept"}
	for {
		infolog.Println("Sending finish")
		res, err := coprocessor.SendFinish(&finish)
		if err != nil {
			errlog.Panicln("Error: error making http request: ", err)
		}
		infolog.Println("Received finish status ", strconv.Itoa(res.StatusCode))

		if res.StatusCode == 202 {
			infolog.Println("No pending rollup request, trying again")
		} else {

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				errlog.Panicln("Error: could not read response body: ", err)
			}

			var finishResponse coprocessor.FinishResponse
			err = json.Unmarshal(resBody, &finishResponse)
			if err != nil {
				errlog.Panicln("Error: unmarshaling body:", err)
			}

			var rawPayload struct {
				Data string `json:"payload"`
			}
			if err := json.Unmarshal(finishResponse.Data, &rawPayload); err != nil {
				errlog.Println("Error unmarshaling payload:", err)
				finish.Status = "reject"
			}

			finish.Status = "accept"
			advanceResponse, err := coprocessor.EvmAdvanceParser(rawPayload.Data)
			if err != nil {
				errlog.Println(err)
				finish.Status = "reject"
			}

			err = Handler(&advanceResponse)
			if err != nil {
				errlog.Println(err)
				finish.Status = "reject"
			}
		}
	}
}
