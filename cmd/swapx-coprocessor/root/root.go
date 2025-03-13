package root

import (
	"encoding/json"
	"io"
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
		slog.Error("Error: could not setup in-memory DB", "err", err)
	}
	slog.Info("In-memory database initialized")

	oh, err := NewMatchOrdersHandler(db, ROLLUP_HTTP_SERVER_URL)
	if err != nil {
		slog.Error("Failed to initialize OrderHandler: %v", "err", err)
	}
	slog.Info("Order handler initialized")

	for {
		slog.Info("Sending finish request")
		finish := coprocessor.FinishRequest{Status: "accept"}
		res, err := coprocessor.SendFinish(&finish)
		if err != nil {
			slog.Error("Error: making HTTP request: %v", "err", err)
		}
		slog.Info("Received finish status", "status_code", strconv.Itoa(res.StatusCode))

		if res.StatusCode == 202 {
			slog.Info("No pending rollup request, retrying...")
			time.Sleep(1 * time.Second)
			continue
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			slog.Error("Error: could not read response body: %v", "err", err)
		}

		var finishResponse coprocessor.FinishResponse
		if err := json.Unmarshal(resBody, &finishResponse); err != nil {
			slog.Error("Error: unmarshaling response body: %v", "err", err)
		}

		var rawPayload struct {
			Data string `json:"payload"`
		}
		if err := json.Unmarshal(finishResponse.Data, &rawPayload); err != nil {
			slog.Error("Error unmarshaling payload", "err", err)
			finish.Status = "reject"
			continue
		}

		advanceResponse, err := coprocessor.EvmAdvanceParser(rawPayload.Data)
		if err != nil {
			slog.Error("Error parsing advance response", "err", err)
			finish.Status = "reject"
			continue
		}

		if err := oh.MatchOrdersHandler(&advanceResponse); err != nil {
			slog.Error("Error handling order book", "err", err)
			finish.Status = "reject"
		}
	}
}
