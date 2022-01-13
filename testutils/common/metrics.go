package common

import (
	"sync"
	"testing"
)

// TestMetricsCallback is the metrics callback to use in tests.
type TestMetricsCallback struct {
	lock *sync.Mutex

	clientBytesReceived  int
	clientBytesSent      int
	clientConnects       int
	clientErrors         int
	messageSendFailures  int
	messageSendSuccesses int
	reconnectAttempts    int
	reconnectFailures    int
	reconnectSuccesses   int
}

// NewTestMetricsCallback creates a new configured instance of test metrics callback.
func NewTestMetricsCallback(t *testing.T) *TestMetricsCallback {
	return &TestMetricsCallback{lock: &sync.Mutex{}}
}

// -- metrics callback interface

func (p *TestMetricsCallback) ClientBytesReceived(n int) {
	p.lock.Lock()
	p.clientBytesReceived = p.clientBytesReceived + n
	p.lock.Unlock()
}
func (p *TestMetricsCallback) ClientBytesSent(n int) {
	p.lock.Lock()
	p.clientBytesSent = p.clientBytesSent + n
	p.lock.Unlock()
}
func (p *TestMetricsCallback) ClientConnect() {
	p.lock.Lock()
	p.clientConnects = p.clientConnects + 1
	p.lock.Unlock()
}
func (p *TestMetricsCallback) ClientError() {
	p.lock.Lock()
	p.clientErrors = p.clientErrors + 1
	p.lock.Unlock()
}
func (p *TestMetricsCallback) ClientMessageSendFailure() {
	p.lock.Lock()
	p.messageSendFailures = p.messageSendFailures + 1
	p.lock.Unlock()
}
func (p *TestMetricsCallback) ClientMessageSendSuccess() {
	p.lock.Lock()
	p.messageSendSuccesses = p.messageSendSuccesses + 1
	p.lock.Unlock()
}
func (p *TestMetricsCallback) ClientReconnectAttempt() {
	p.lock.Lock()
	p.reconnectAttempts = p.reconnectAttempts + 1
	p.lock.Unlock()
}
func (p *TestMetricsCallback) ClientReconnectFailure() {
	p.lock.Lock()
	p.reconnectFailures = p.reconnectFailures + 1
	p.lock.Unlock()
}
func (p *TestMetricsCallback) ClientReconnectSuccess() {
	p.lock.Lock()
	p.reconnectSuccesses = p.reconnectSuccesses + 1
	p.lock.Unlock()
}

// -- metrics inspect interface

func (p *TestMetricsCallback) InspectClientBytesReceived(t *testing.T) int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.clientBytesReceived
}
func (p *TestMetricsCallback) InspectClientBytesSent(t *testing.T) int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.clientBytesSent
}
func (p *TestMetricsCallback) InspectClientConnect(t *testing.T) int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.clientConnects
}
func (p *TestMetricsCallback) InspectClientError(t *testing.T) int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.clientErrors
}
func (p *TestMetricsCallback) InspectClientMessageSendFailure(t *testing.T) int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.messageSendFailures
}
func (p *TestMetricsCallback) InspectClientMessageSendSuccess(t *testing.T) int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.messageSendSuccesses
}
func (p *TestMetricsCallback) InspectClientReconnectAttempt(t *testing.T) int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.reconnectAttempts
}
func (p *TestMetricsCallback) InspectClientReconnectFailure(t *testing.T) int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.reconnectFailures
}
func (p *TestMetricsCallback) InspectClientReconnectSuccess(t *testing.T) int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.reconnectSuccesses
}
