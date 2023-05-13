package dns

import (
	"crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"
	"net"
	"strings"
)

const rootNS = "198.41.0.4"

type Resolver struct {
	history map[string]bool
	indent  int
	domain  string
	ns      string
}

func NewResolver(domain string) *Resolver {
	return &Resolver{
		history: make(map[string]bool),
		domain:  domain,
		ns:      rootNS,
	}
}

func (r *Resolver) Resolve() (net.IP, error) {
	r.indent++

	for {
		msg, err := r.doResolve()
		if err != nil {
			return nil, err
		}

		ip, err := r.handle(msg)
		if err != nil {
			return nil, err
		}

		if ip != nil {
			r.indent--

			return ip, nil
		}
	}
}

func (r *Resolver) doResolve() (*message, error) {
	r.log("Resolving %s on %s...", r.domain, r.ns)

	key := fmt.Sprintf("%s@%s", r.domain, r.ns)
	if r.history[key] {
		return nil, NewError("endless loop detected")
	}

	r.history[key] = true

	q, err := r.query()
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("udp", r.ns+":53")
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	if err = q.sendRequest(conn); err != nil {
		return nil, err
	}

	msg, err := q.readResponse(conn)

	return msg, err
}

func (r *Resolver) handle(msg *message) (net.IP, error) {
	answer, err := msg.answer()
	if err != nil {
		return nil, err
	}

	switch answer.recType { //nolint:exhaustive
	case A:
		return answer.dataAsIP(), nil
	case CNAME:
		r.domain = answer.dataAsString()
		r.ns = rootNS

		r.log("=> should query %s", r.domain)

		return nil, nil
	case NS:
		ip, err := msg.firstIP()
		if err == nil {
			r.ns = ip.String()

			r.log("=> should ask %s", r.ns)

			return nil, nil
		}

		ns := answer.dataAsString()

		r.log(err.Error())
		r.log("=> should ask %s", ns)

		d := r.domain
		r.domain = ns
		r.ns = rootNS

		ip, err = r.Resolve()
		if err != nil {
			return nil, err
		}

		r.domain = d
		r.ns = ip.String()

		r.log("=> %s = %s", ns, r.ns)

		return nil, nil
	}

	return nil, NewError("don't know how to handle answer %#v", msg)
}

func (r *Resolver) query() (*query, error) {
	rnd, err := rand.Int(rand.Reader, big.NewInt(math.MaxUint16+1))
	if err != nil {
		return nil, err
	}

	result := &query{
		h: &header{
			ID:           uint16(rnd.Uint64()),
			NumQuestions: 1,
		},
		q: &question{
			name:    r.domain,
			recType: A,
			class:   IN,
		},
	}

	return result, nil
}

func (r *Resolver) log(format string, a ...any) {
	log.Printf(strings.Repeat("    ", r.indent-1)+format, a...)
}
