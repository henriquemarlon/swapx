package root

import (
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/henriquemarlon/swapx/configs"
	"github.com/henriquemarlon/swapx/pkg/coprocessor"
	"github.com/spf13/cobra"
)

const (
	CMD_NAME = "swapx-coprocessor"
)

var (
	verbose bool
	Cmd     = &cobra.Command{
		Use:   CMD_NAME,
		Short: "Run SwapX Coprocessor",
		Long:  `EVM Linux Coprocessor as an orderbook for UniswapV4 Hooks`,
		Run:   run,
	}
	ROLLUP_HTTP_SERVER_URL = os.Getenv("ROLLUP_HTTP_SERVER_URL")
)

func init() {
	Cmd.Flags().BoolVar(&verbose, "verbose", false, "Show detailed logs")
	Cmd.PreRun = func(cmd *cobra.Command, args []string) {
		configs.ConfigureLogger(slog.LevelInfo)
	}
}

func run(cmd *cobra.Command, args []string) {
	db, err := configs.SetupInMemoryDB()
	if err != nil {
		log.Fatalf("Error: could not setup in-memory DB: %v", err)
	}
	log.Println("In-memory database initialized")

	oh, err := NewOrderBookHandler(db, ROLLUP_HTTP_SERVER_URL)
	if err != nil {
		log.Fatalf("Failed to initialize OrderHandler: %v", err)
	}
	log.Println("Order handler initialized")

	for {
		log.Println("Sending finish request")
		finish := coprocessor.FinishRequest{Status: "accept"}
		res, err := coprocessor.SendFinish(&finish)
		if err != nil {
			log.Fatalf("Error: making HTTP request: %v", err)
		}
		log.Println("Received finish status", strconv.Itoa(res.StatusCode))

		if res.StatusCode == 202 {
			log.Println("No pending rollup request, retrying...")
			time.Sleep(1 * time.Second)
			continue
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("Error: could not read response body: %v", err)
		}

		var finishResponse coprocessor.FinishResponse
		if err := json.Unmarshal(resBody, &finishResponse); err != nil {
			log.Fatalf("Error: unmarshaling response body: %v", err)
		}

		var rawPayload struct {
			Data string `json:"payload"`
		}
		if err := json.Unmarshal(finishResponse.Data, &rawPayload); err != nil {
			log.Println("Error unmarshaling payload", err)
			finish.Status = "reject"
			continue
		}

		advanceResponse, err := coprocessor.EvmAdvanceParser(rawPayload.Data)
		if err != nil {
			log.Println("Error parsing advance response", err)
			finish.Status = "reject"
			continue
		}

		if err := oh.OrderBookHandler(&advanceResponse); err != nil {
			log.Println("Error handling order", err)
			finish.Status = "reject"
		}
	}
}
