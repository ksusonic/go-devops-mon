package trust

import (
	"net"
	"net/http"

	"go.uber.org/zap"
)

type NetTrustService struct {
	trustedAddress net.IPNet
	logger         *zap.Logger
}

func NewNetTrustService(trustedAddress string, logger *zap.Logger) (*NetTrustService, error) {
	_, mask, err := net.ParseCIDR(trustedAddress)
	if err != nil {
		return nil, err
	}
	return &NetTrustService{
		trustedAddress: *mask,
		logger:         logger,
	}, nil
}

func (s NetTrustService) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RemoteAddr == "" {
			http.Error(w, "empty ip is forbidden", http.StatusForbidden)
			return
		}

		if !s.trustedAddress.Contains(net.ParseIP(r.RemoteAddr)) {
			s.logger.Error("got forbidden for ip", zap.String("remote_ip", r.RemoteAddr), zap.String("mask", s.trustedAddress.String()))
			http.Error(w, "ip subnet is forbidden", http.StatusForbidden)
		}
		next.ServeHTTP(w, r)
	})
}
