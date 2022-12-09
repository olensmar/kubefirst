package ngrok

import (
	"context"
	"github.com/rs/zerolog/log"
	"io"
	"net"

	"github.com/spf13/viper"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
	"golang.org/x/sync/errgroup"
)

type Tunnel struct {
	ctx           context.Context
	cancelContext context.CancelFunc
	localDestiny  string
}

// NewNgrok create a Ngrok object
func NewNgrok(ctx context.Context, cancelContext context.CancelFunc) *Tunnel {
	return &Tunnel{
		ctx:           ctx,
		cancelContext: cancelContext,
	}
}

// OpenTunnel opens a TCP connection on a randomized Ngrok DNS, and forward request to the destiny(dest) local address.
func (t Tunnel) OpenTunnel() {

	tunnel, err := ngrok.Listen(t.ctx, config.HTTPEndpoint(), ngrok.WithAuthtokenFromEnv())
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	log.Info().Msgf("tunnel created: %s", tunnel.URL())
	viper.Set("github.atlantis.webhook.url", tunnel.URL()+"/events")
	viper.Set("ngrok.url", tunnel.URL())
	if err := viper.WriteConfig(); err != nil {
		log.Error().Err(err).Msg("")
	}

	conn, err := tunnel.Accept()
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	log.Info().Msgf("Ngrok tunnel is open and accepting new requests %v...", conn.RemoteAddr())

	go func() {
		select {
		case <-t.ctx.Done():
			if err := tunnel.CloseWithContext(t.ctx); err != nil {
				log.Error().Err(err).Msg("")
				return
			}
			log.Info().Msg("Ngrok tunnel is closed, and not accepting new requests")
		}
	}()
	err = handleConn(t.ctx, conn)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
}

// handleConn accepts connection request, and forward Ngrok tunnel to internal localDestiny.
func handleConn(ctx context.Context, conn net.Conn) error {
	next, err := net.Dial("tcp", ":80")
	if err != nil {
		return err
	}

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		_, err := io.Copy(next, conn)
		return err
	})
	g.Go(func() error {
		_, err := io.Copy(conn, next)
		return err
	})

	return g.Wait()
}
