package ogmigo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	//"os"
	//"sort"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

type MonitorMempool struct {
	cancel context.CancelFunc
	errs   chan error
	done   chan struct{}
	err    error
	logger Logger
}

func (c *MonitorMempool) Done() <-chan struct{} {
	return c.done
}

func (c *MonitorMempool) Err() <-chan error {
	return c.errs
}

func (c *MonitorMempool) Close() error {
	c.cancel()
	select {
	case v := <-c.errs:
		c.err = v
	default:
		// err already set
	}
	return c.err
}

type MonitorMempoolFunc func(ctx context.Context, data []byte) (bool, error)

type MonitorMempoolOptions struct {
	reconnect bool // reconnect to ogmios if connection drops
}

func buildMonitorMempoolOptions(opts ...MonitorMempoolOption) MonitorMempoolOptions {
	var options MonitorMempoolOptions
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

type MonitorMempoolOption func(opts *MonitorMempoolOptions)

func (c *Client) MonitorMempool(ctx context.Context, callback MonitorMempoolFunc, opts ...MonitorMempoolOption) (*MonitorMempool, error) {
	options := buildMonitorMempoolOptions(opts...)

	done := make(chan struct{})
	errs := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer close(done)

		var (
			timeout = 10 * time.Second
			err     error
		)
		for {
			err = c.doMonitorMempool(ctx, callback, options)
			if err != nil && isTemporaryError(err) {
				if options.reconnect {
					c.options.logger.Info("websocket connection error: will retry",
						KV("delay", timeout.Round(time.Millisecond).String()),
						KV("err", err.Error()),
					)

					select {
					case <-ctx.Done():
						return
					case <-time.After(timeout):
						continue
					}
				}
			}

			break
		}
		errs <- err
	}()

	return &MonitorMempool{
		cancel: cancel,
		errs:   errs,
		done:   done,
		logger: c.logger,
	}, nil
}

type NextTransaction struct {
        Result *chainsync.Tx
}

func (c *Client) doMonitorMempool(ctx context.Context, callback MonitorMempoolFunc, options MonitorMempoolOptions) error {
	conn, _, err := websocket.DefaultDialer.Dial(c.options.endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to ogmios, %v: %w", c.options.endpoint, err)
	}

	init, err := json.Marshal(Map{
		"jsonrpc": "2.0",
		"method":  "acquireMempool",
		"id":      Map{"step": "MEMPOOLINIT"},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal init message: %w", err)
	}

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		c.options.logger.Info("ogmigo mempool monitoring started")
		defer c.options.logger.Info("ogmigo mempool monitoring stopped")
		<-ctx.Done()
		return nil
	})

	var connState int64 // 0 - open, 1 - closing, 2 - closed
	group.Go(func() error {
		<-ctx.Done()
		atomic.AddInt64(&connState, 1)
		if err := conn.Close(); err != nil {
			return err
		}
		atomic.AddInt64(&connState, 1)
		return nil
	})

	// prime the pump
	ch := make(chan struct{}, 64)
	for i := 0; i < c.options.pipeline; i++ {
		select {
		case ch <- struct{}{}:
		default:
		}
	}

	group.Go(func() error {
		if err := conn.WriteMessage(websocket.TextMessage, init); err != nil {
			var oe *net.OpError
			if ok := errors.As(err, &oe); ok {
				if v := atomic.LoadInt64(&connState); v > 0 {
					return nil // connection closed
				}
			}
			return fmt.Errorf("failed to write AcquireMempool: %w", err)
		}

		next := []byte(`{"jsonrpc":"2.0","method":"nextTransaction","params":{"fields":"all"},"id":{}}`)
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ch:
				if err := conn.WriteMessage(websocket.TextMessage, next); err != nil {
					return fmt.Errorf("failed to write nextTransaction: %w", err)
				}
			}
		}
	})

	group.Go(func() error {
		for n := uint64(1); ; n++ {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				var oe *net.OpError
				if ok := errors.As(err, &oe); ok {
					if v := atomic.LoadInt64(&connState); v > 0 {
						return nil // connection closed
					}
				}
				return fmt.Errorf("failed to read message from ogmios: %w", err)
			}

			switch messageType {
			case websocket.BinaryMessage:
				c.options.logger.Info("skipping unexpected binary message")
				continue

			case websocket.CloseMessage:
				return nil

			case websocket.PingMessage:
				if err := conn.WriteMessage(websocket.PongMessage, nil); err != nil {
					return fmt.Errorf("failed to respond with pong to ogmios: %w", err)
				}
				continue

			case websocket.PongMessage:
				continue

			case websocket.TextMessage:
				// ok
			}

                        //var nextTransactionResponse Response
                        //if err := json.Unmarshal(data, &nextTransactionResponse); err != nil {
                        //        return fmt.Errorf("couldn't unmarshal response into Response: %w", err)
                        //}

                        //if nextTransactionResponse.Result == nil {
                        //       // Acquire mempool again 
                        //} else {
                        //       err := callback(ctx, nextTransactionResponse.Result)
                        //       if err != nil {
                        //               return fmt.Errorf("mempool monitoring stopped: callback failed: %w", err)
                        //       }
                        //}

			requestNext, err := callback(ctx, data)
			if err != nil {
				return fmt.Errorf("mempool monitoring stopped: callback failed: %w", err)
			}
			if requestNext {
				select {
				case <-ctx.Done():
					return nil
				case ch <- struct{}{}:
					// request the next message
				default:
					// pump is full
				}
			} else {
				return nil
			}
		}
	})
	return group.Wait()
}
