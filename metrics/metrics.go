package metrics

// Callback represents a metrics callback
// interface used internally by the client.
// Provide your own implementation to retrieve metrics
// produced by the client.
type Callback interface {
	ClientBytesReceived(n int)
	ClientBytesSent(n int)
	ClientConnect()
	ClientError()
	ClientMessageSendFailure()
	ClientMessageSendSuccess()
	ClientReconnectAttempt()
	ClientReconnectFailure()
	ClientReconnectSuccess()
}

// Noop returns an instance of noop metric
func Noop() Callback {
	return &noop{}
}

type noop struct {
}

func (p *noop) ClientBytesReceived(n int) {}
func (p *noop) ClientBytesSent(n int)     {}
func (p *noop) ClientConnect()            {}
func (p *noop) ClientError()              {}
func (p *noop) ClientMessageSendFailure() {}
func (p *noop) ClientMessageSendSuccess() {}
func (p *noop) ClientReconnectAttempt()   {}
func (p *noop) ClientReconnectFailure()   {}
func (p *noop) ClientReconnectSuccess()   {}
