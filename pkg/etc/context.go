package etc

import (
	"context"
	"fmt"
	"time"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
	"google.golang.org/grpc/metadata"
)

func NewTimoutContext(ctx context.Context) (context.Context, context.CancelFunc) {

	md := metadata.Pairs()

	for _, key := range []string{config.KeyUserId, config.KeyCompanyId, config.KeyTimezone} {
		if ctx.Value(key) != nil {
			fmt.Println("key_timezone", ctx.Value(key))
			val, ok := ctx.Value(key).(string)
			if ok {
				md.Set(key, val)
			}
		}
	}

	ctx = metadata.NewOutgoingContext(ctx, md)

	res, cancel := context.WithTimeout(ctx, time.Second*config.GrpcContextTimeout)

	return res, cancel
}
