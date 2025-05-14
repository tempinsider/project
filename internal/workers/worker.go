package workers

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Worker struct {
	ID           uuid.UUID
	Interval     time.Duration
	mutex        sync.Mutex
	ctx          context.Context
	counter      float64
	ticker       *time.Ticker
	closer       chan bool
	status       bool
	workFunc     func(conn *pgx.Conn)
	dbConnection *pgx.Conn
}

func Create(ctx context.Context, interval time.Duration, workFunc func(conn *pgx.Conn), conn *pgx.Conn) (*Worker, error) {
	worker := &Worker{ctx: ctx, ID: uuid.New(), Interval: interval, closer: make(chan bool), workFunc: workFunc, dbConnection: conn}
	return worker, nil
}

func (w *Worker) Start() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.ticker = time.NewTicker(time.Second)
	w.counter = w.Interval.Seconds()
	w.scheduler()
	w.status = true
}

func (w *Worker) Stop() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.ticker.Stop()
	w.closer <- true
	w.counter = 0
	w.status = false
}

func (w *Worker) scheduler() {
	targetDuration := w.Interval.Seconds()
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				return
			case <-w.closer: // ticker doesn't provide a channel after we stop it, so we simply had to do it manually.
				return
			case <-w.ticker.C:
				w.counter++
				if w.counter >= targetDuration {
					w.work()
					w.counter = 0
				}
			}
		}
	}()
}

func (w *Worker) work() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.workFunc(w.dbConnection)
}

func (w *Worker) IsWorking() bool {
	return w.status
}
