package utils

import (
	"sync"
	"time"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

type HandlerFunc func([]byte)

type PubSub struct {
	lsn  *pq.Listener
	log  *zap.Logger
	fns  map[string][]HandlerFunc
	mu   sync.RWMutex
	done chan struct{}
}

func NewPubSub(connString string, logger *zap.Logger) *PubSub {
	listener := pq.NewListener(connString,
		10*time.Second,
		time.Minute,
		func(ev pq.ListenerEventType, err error) {
			if err != nil {
				logger.Error("listener error", zap.Error(err))
			}
		})

	return &PubSub{
		lsn:  listener,
		log:  logger,
		fns:  make(map[string][]HandlerFunc),
		done: make(chan struct{}),
	}
}

func (p *PubSub) Subscribe(topic string, handler HandlerFunc) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.fns[topic]) == 0 {
		err := p.lsn.Listen(topic)
		if err != nil {
			p.log.Error("failed to listen to topic",
				zap.String("topic", topic),
				zap.Error(err))
			return
		}
	}

	p.fns[topic] = append(p.fns[topic], handler)
}

func (p *PubSub) Start() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-p.done:
				return
			case <-ticker.C:
				err := p.lsn.Ping()
				if err != nil {
					p.log.Error("ping failed", zap.Error(err))
				}
			case n := <-p.lsn.Notify:
				if n == nil {
					continue
				}

				p.log.Debug("notification",
					zap.String("topic", n.Channel),
					zap.String("payload", n.Extra),
				)

				p.mu.RLock()
				handlers := p.fns[n.Channel]
				p.mu.RUnlock()

				for _, handler := range handlers {
					go func(fn HandlerFunc, payload []byte) {
						defer func() {
							if r := recover(); r != nil {
								p.log.Error("handler panic",
									zap.Any("recover", r),
									zap.String("channel", n.Channel))
							}
						}()
						fn(payload)
					}(handler, []byte(n.Extra))
				}
			}
		}
	}()
}

func (p *PubSub) Stop() {
	close(p.done)
	err := p.lsn.Close()
	if err != nil {
		p.log.Error("error closing listener", zap.Error(err))
	}
}
