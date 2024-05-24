package main

import (
	//"bytes"
	"context"
	//"encoding/json"
	//"flag"
	"fmt"
	"os"
	"time"

	"github.com/SundaeSwap-finance/ogmigo/v6"
)

func main() {
	n := 10
	var callback ogmigo.MonitorMempoolFunc = func(ctx context.Context, data []byte) (bool, error) {
		fmt.Printf("%s\n", string(data))
		if n > 0 {
                        n = n - 1
			return true, nil
		} else {
			return false, nil
		}
	}

	//debugPtr := flag.Bool("debug", false, "Print debug statements")
	//flag.Parse()
	ctx := context.Background()

	ogmiosAddr := os.Getenv("OGMIOS")
	client := ogmigo.New(ogmigo.WithEndpoint(ogmiosAddr))
	closer, err := client.MonitorMempool(ctx, callback)
	if err != nil {
		fmt.Printf("Failed MonitorMempool Open: %v\n", err)
		return
	}

	time.Sleep(2 * time.Second)

	if err := closer.Close(); err != nil {
		fmt.Printf("Failed MonitorMempool Close: %v\n", err)
	}
}
