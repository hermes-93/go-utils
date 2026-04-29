package health

import (
	"fmt"
	"net/http"
	"time"
)

type Options struct {
	Timeout       int
	Retries       int
	RetryInterval int
	ExpectedCode  int
}

type Result struct {
	URL        string
	Healthy    bool
	StatusCode int
	Latency    time.Duration
	Error      string
	Attempts   int
}

func (r Result) String() string {
	if r.Healthy {
		return fmt.Sprintf("OK  %s — HTTP %d — %dms (attempt %d/%d)",
			r.URL, r.StatusCode, r.Latency.Milliseconds(), r.Attempts, r.Attempts)
	}
	if r.Error != "" {
		return fmt.Sprintf("ERR %s — %s (after %d attempt(s))", r.URL, r.Error, r.Attempts)
	}
	return fmt.Sprintf("ERR %s — HTTP %d, expected %d (after %d attempt(s))",
		r.URL, r.StatusCode, r.Attempts, r.Attempts)
}

func Check(url string, opts Options) Result {
	client := &http.Client{
		Timeout: time.Duration(opts.Timeout) * time.Second,
	}

	maxAttempts := opts.Retries + 1
	var lastErr string
	var lastCode int
	var lastLatency time.Duration

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		start := time.Now()
		resp, err := client.Get(url)
		latency := time.Since(start)

		if err != nil {
			lastErr = err.Error()
			lastLatency = latency
			if attempt < maxAttempts && opts.RetryInterval > 0 {
				time.Sleep(time.Duration(opts.RetryInterval) * time.Second)
			}
			continue
		}
		resp.Body.Close()

		lastCode = resp.StatusCode
		lastLatency = latency
		lastErr = ""

		if resp.StatusCode == opts.ExpectedCode {
			return Result{
				URL:        url,
				Healthy:    true,
				StatusCode: resp.StatusCode,
				Latency:    latency,
				Attempts:   attempt,
			}
		}
		if attempt < maxAttempts && opts.RetryInterval > 0 {
			time.Sleep(time.Duration(opts.RetryInterval) * time.Second)
		}
	}

	return Result{
		URL:        url,
		Healthy:    false,
		StatusCode: lastCode,
		Latency:    lastLatency,
		Error:      lastErr,
		Attempts:   maxAttempts,
	}
}
