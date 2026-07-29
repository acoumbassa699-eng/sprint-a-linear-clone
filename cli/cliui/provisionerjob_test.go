package cliui_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/cliui"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtime"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil/expecter"
	"github.com/optimus-ide-collab/serpent"
)

// This cannot be ran in parallel because it uses a signal.
// nolint:tparallel
func TestProvisionerJob(t *testing.T) {
	t.Run("NoLogs", func(t *testing.T) {
		t.Parallel()

		test := newProvisionerJob(t)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		testutil.Go(t, func() {
			<-test.Next
			test.JobMutex.Lock()
			test.Job.Status = optimus-ide-collabsdk.ProvisionerJobRunning
			now := dbtime.Now()
			test.Job.StartedAt = &now
			test.JobMutex.Unlock()
			<-test.Next
			test.JobMutex.Lock()
			test.Job.Status = optimus-ide-collabsdk.ProvisionerJobSucceeded
			now = dbtime.Now()
			test.Job.CompletedAt = &now
			close(test.Logs)
			test.JobMutex.Unlock()
		})
		testutil.Eventually(ctx, t, func(ctx context.Context) (done bool) {
			test.Stdout.ExpectMatch(ctx, cliui.ProvisioningStateQueued)
			test.Next <- struct{}{}
			test.Stdout.ExpectMatch(ctx, cliui.ProvisioningStateQueued)
			test.Stdout.ExpectMatch(ctx, cliui.ProvisioningStateRunning)
			test.Next <- struct{}{}
			test.Stdout.ExpectMatch(ctx, cliui.ProvisioningStateRunning)
			return true
		}, testutil.IntervalFast)
	})

	t.Run("Stages", func(t *testing.T) {
		t.Parallel()

		test := newProvisionerJob(t)
		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		testutil.Go(t, func() {
			<-test.Next
			test.JobMutex.Lock()
			test.Job.Status = optimus-ide-collabsdk.ProvisionerJobRunning
			now := dbtime.Now()
			test.Job.StartedAt = &now
			test.Logs <- optimus-ide-collabsdk.ProvisionerJobLog{
				CreatedAt: dbtime.Now(),
				Stage:     "Something",
			}
			test.JobMutex.Unlock()
			<-test.Next
			test.JobMutex.Lock()
			test.Job.Status = optimus-ide-collabsdk.ProvisionerJobSucceeded
			now = dbtime.Now()
			test.Job.CompletedAt = &now
			close(test.Logs)
			test.JobMutex.Unlock()
		})
		testutil.Eventually(ctx, t, func(ctx context.Context) (done bool) {
			test.Stdout.ExpectMatch(ctx, cliui.ProvisioningStateQueued)
			test.Next <- struct{}{}
			test.Stdout.ExpectMatch(ctx, cliui.ProvisioningStateQueued)
			test.Stdout.ExpectMatch(ctx, "Something")
			test.Next <- struct{}{}
			test.Stdout.ExpectMatch(ctx, "Something")
			return true
		}, testutil.IntervalFast)
	})

	t.Run("Queue Position", func(t *testing.T) {
		t.Parallel()

		stage := cliui.ProvisioningStateQueued

		tests := []struct {
			name     string
			queuePos int
			expected string
		}{
			{
				name:     "first",
				queuePos: 0,
				expected: fmt.Sprintf("%s$", stage),
			},
			{
				name:     "next",
				queuePos: 1,
				expected: fmt.Sprintf(`%s %s$`, stage, regexp.QuoteMeta("(next)")),
			},
			{
				name:     "other",
				queuePos: 4,
				expected: fmt.Sprintf(`%s %s$`, stage, regexp.QuoteMeta("(position: 4)")),
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				test := newProvisionerJob(t)
				test.JobMutex.Lock()
				test.Job.QueuePosition = tc.queuePos
				test.Job.QueueSize = tc.queuePos
				test.JobMutex.Unlock()

				ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
				defer cancel()

				testutil.Go(t, func() {
					<-test.Next
					test.JobMutex.Lock()
					test.Job.Status = optimus-ide-collabsdk.ProvisionerJobRunning
					now := dbtime.Now()
					test.Job.StartedAt = &now
					test.JobMutex.Unlock()
					<-test.Next
					test.JobMutex.Lock()
					test.Job.Status = optimus-ide-collabsdk.ProvisionerJobSucceeded
					now = dbtime.Now()
					test.Job.CompletedAt = &now
					close(test.Logs)
					test.JobMutex.Unlock()
				})
				testutil.Eventually(ctx, t, func(ctx context.Context) (done bool) {
					test.Stdout.ExpectRegexMatch(ctx, tc.expected)
					test.Next <- struct{}{}
					test.Stdout.ExpectMatch(ctx, cliui.ProvisioningStateQueued) // step completed
					test.Stdout.ExpectMatch(ctx, cliui.ProvisioningStateRunning)
					test.Next <- struct{}{}
					test.Stdout.ExpectMatch(ctx, cliui.ProvisioningStateRunning)
					return true
				}, testutil.IntervalFast)
			})
		}
	})

	// This cannot be ran in parallel because it uses a signal.
	// nolint:paralleltest
	t.Run("Cancel", func(t *testing.T) {
		t.Skip("This test issues an interrupt signal which will propagate to the test runner.")

		if runtime.GOOS == "windows" {
			// Sending interrupt signal isn't supported on Windows!
			t.SkipNow()
		}

		test := newProvisionerJob(t)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		testutil.Go(t, func() {
			<-test.Next
			currentProcess, err := os.FindProcess(os.Getpid())
			assert.NoError(t, err)
			err = currentProcess.Signal(os.Interrupt)
			assert.NoError(t, err)
			<-test.Next
			test.JobMutex.Lock()
			test.Job.Status = optimus-ide-collabsdk.ProvisionerJobCanceled
			now := dbtime.Now()
			test.Job.CompletedAt = &now
			close(test.Logs)
			test.JobMutex.Unlock()
		})
		testutil.Eventually(ctx, t, func(ctx context.Context) (done bool) {
			test.Stdout.ExpectMatch(ctx, cliui.ProvisioningStateQueued)
			test.Next <- struct{}{}
			test.Stdout.ExpectMatch(ctx, "Gracefully canceling")
			test.Next <- struct{}{}
			test.Stdout.ExpectMatch(ctx, cliui.ProvisioningStateQueued)
			return true
		}, testutil.IntervalFast)
	})
}

type provisionerJobTest struct {
	Next     chan struct{}
	Job      *optimus-ide-collabsdk.ProvisionerJob
	JobMutex *sync.Mutex
	Logs     chan optimus-ide-collabsdk.ProvisionerJobLog
	Stdout   *expecter.Expecter
}

func newProvisionerJob(t *testing.T) provisionerJobTest {
	job := &optimus-ide-collabsdk.ProvisionerJob{
		Status:    optimus-ide-collabsdk.ProvisionerJobPending,
		CreatedAt: dbtime.Now(),
	}
	jobLock := sync.Mutex{}
	logs := make(chan optimus-ide-collabsdk.ProvisionerJobLog, 1)
	cmd := &serpent.Command{
		Handler: func(inv *serpent.Invocation) error {
			return cliui.ProvisionerJob(inv.Context(), inv.Stdout, cliui.ProvisionerJobOptions{
				FetchInterval: time.Millisecond,
				Fetch: func() (optimus-ide-collabsdk.ProvisionerJob, error) {
					jobLock.Lock()
					defer jobLock.Unlock()
					return *job, nil
				},
				Cancel: func() error {
					return nil
				},
				Logs: func() (<-chan optimus-ide-collabsdk.ProvisionerJobLog, io.Closer, error) {
					return logs, closeFunc(func() error {
						return nil
					}), nil
				},
			})
		},
	}
	inv := cmd.Invoke()

	stdout := expecter.NewAttachedToInvocation(t, inv)
	done := make(chan struct{})
	go func() {
		defer close(done)
		err := inv.WithContext(context.Background()).Run()
		if err != nil {
			assert.ErrorIs(t, err, cliui.ErrCanceled)
		}
	}()
	t.Cleanup(func() {
		<-done
	})
	return provisionerJobTest{
		Next:     make(chan struct{}),
		Job:      job,
		JobMutex: &jobLock,
		Logs:     logs,
		Stdout:   stdout,
	}
}

type closeFunc func() error

func (c closeFunc) Close() error {
	return c()
}
