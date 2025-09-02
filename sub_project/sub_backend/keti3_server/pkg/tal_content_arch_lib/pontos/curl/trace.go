package curl

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
	"net/http/httptrace"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
)

type timeString time.Duration

func (d timeString) Value() time.Duration {
	return time.Duration(d)
}

func (d timeString) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *timeString) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = timeString(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = timeString(tmp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// traceInfo struct
//_______________________________________________________________________

// traceInfo struct is used provide request trace info such as DNS lookup
// duration, Connection obtain duration, Server processing duration, etc.
type traceInfo struct {
	// DNSLookup is a duration that transport took to perform
	// DNS lookup.
	DNSLookup timeString `json:"dns_lookup,omitempty"`

	// ConnTime is a duration that took to obtain a successful connection.
	ConnTime timeString `json:"conn_time,omitempty"`

	// TCPConnTime is a duration that took to obtain the TCP connection.
	TCPConnTime timeString `json:"tcp_conn_time,omitempty"`

	// TLSHandshake is a duration that TLS handshake took place.
	TLSHandshake timeString `json:"tls_handshake,omitempty"`

	// ServerTime is a duration that server took to respond first byte.
	ServerTime timeString `json:"server_time,omitempty"`

	// ResponseTime is a duration since first response byte from server to
	// request completion.
	ResponseTime timeString `json:"response_time,omitempty"`

	// TotalTime is a duration that total request took end-to-end.
	TotalTime timeString `json:"total_time,omitempty"`

	// IsConnReused is whether this connection has been previously
	// used for another HTTP request.
	IsConnReused bool `json:"is_conn_reused,omitempty"`

	// IsConnWasIdle is whether this connection was obtained from an
	// idle pool.
	IsConnWasIdle bool `json:"is_conn_was_idle,omitempty"`

	// ConnIdleTime is a duration how long the connection was previously
	// idle, if IsConnWasIdle is true.
	ConnIdleTime timeString `json:"conn_idle_time,omitempty"`

	// RemoteAddr returns the remote network address.
	RemoteAddr net.Addr `json:"remote_addr,omitempty"`
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// clientTrace struct and its methods
//_______________________________________________________________________
type clientTrace struct {
	getConn              time.Time
	dnsStart             time.Time
	dnsDone              time.Time
	connectDone          time.Time
	tlsHandshakeStart    time.Time
	tlsHandshakeDone     time.Time
	gotConn              time.Time
	gotFirstResponseByte time.Time
	endTime              time.Time
	gotConnInfo          httptrace.GotConnInfo
}

func (t *clientTrace) createContext(ctx context.Context) context.Context {
	return httptrace.WithClientTrace(
		ctx,
		&httptrace.ClientTrace{
			DNSStart: func(_ httptrace.DNSStartInfo) {
				t.dnsStart = time.Now()
			},
			DNSDone: func(_ httptrace.DNSDoneInfo) {
				t.dnsDone = time.Now()
			},
			ConnectStart: func(_, _ string) {
				if t.dnsDone.IsZero() {
					t.dnsDone = time.Now()
				}
				if t.dnsStart.IsZero() {
					t.dnsStart = t.dnsDone
				}
			},
			ConnectDone: func(net, addr string, err error) {
				t.connectDone = time.Now()
			},
			GetConn: func(_ string) {
				t.getConn = time.Now()
			},
			GotConn: func(ci httptrace.GotConnInfo) {
				t.gotConn = time.Now()
				t.gotConnInfo = ci
			},
			GotFirstResponseByte: func() {
				t.gotFirstResponseByte = time.Now()
			},
			TLSHandshakeStart: func() {
				t.tlsHandshakeStart = time.Now()
			},
			TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
				t.tlsHandshakeDone = time.Now()
			},
		},
	)
}

// traceInfo method returns the trace info for the request.
func (t *clientTrace) traceInfo() traceInfo {
	if t == nil {
		return traceInfo{}
	}

	ti := traceInfo{
		DNSLookup:     timeString(t.dnsDone.Sub(t.dnsStart)),
		TLSHandshake:  timeString(t.tlsHandshakeDone.Sub(t.tlsHandshakeStart)),
		ServerTime:    timeString(t.gotFirstResponseByte.Sub(t.gotConn)),
		IsConnReused:  t.gotConnInfo.Reused,
		IsConnWasIdle: t.gotConnInfo.WasIdle,
		ConnIdleTime:  timeString(t.gotConnInfo.IdleTime),
	}

	// Calculate the total time accordingly,
	// when connection is reused
	if t.gotConnInfo.Reused {
		ti.TotalTime = timeString(t.endTime.Sub(t.getConn))
	} else {
		ti.TotalTime = timeString(t.endTime.Sub(t.dnsStart))
	}

	// Only calculate on successful connections
	if !t.connectDone.IsZero() {
		ti.TCPConnTime = timeString(t.connectDone.Sub(t.dnsDone))
	}

	// Only calculate on successful connections
	if !t.gotConn.IsZero() {
		ti.ConnTime = timeString(t.gotConn.Sub(t.getConn))
	}

	// Only calculate on successful connections
	if !t.gotFirstResponseByte.IsZero() {
		ti.ResponseTime = timeString(t.endTime.Sub(t.gotFirstResponseByte))
	}

	// Capture remote address info when connection is non-nil
	if t.gotConnInfo.Conn != nil {
		ti.RemoteAddr = t.gotConnInfo.Conn.RemoteAddr()
	}

	return ti
}

func newOtelClientTrace(ctx context.Context) context.Context {
	return httptrace.WithClientTrace(
		ctx,
		otelhttptrace.NewClientTrace(
			ctx,
			otelhttptrace.WithoutHeaders(),
		),
	)
}
