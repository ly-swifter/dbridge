package common

import (
	"context"

	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/google/uuid"
	"github.com/lyswifter/dbridge/node/modules/dtypes"
	"go.uber.org/fx"
	"golang.org/x/xerrors"
)

var session = uuid.New()

type CommonAPI struct {
	fx.In

	APISecret    *dtypes.APIAlg
	ShutdownChan dtypes.ShutdownChan
}

type jwtPayload struct {
	Allow []auth.Permission
}

func (a *CommonAPI) AuthVerify(ctx context.Context, token string) ([]auth.Permission, error) {
	var payload jwtPayload
	if _, err := jwt.Verify([]byte(token), (*jwt.HMACSHA)(a.APISecret), &payload); err != nil {
		return nil, xerrors.Errorf("JWT Verification failed: %w", err)
	}

	return payload.Allow, nil
}

func (a *CommonAPI) AuthNew(ctx context.Context, perms []auth.Permission) ([]byte, error) {
	p := jwtPayload{
		Allow: perms, // TODO: consider checking validity
	}

	return jwt.Sign(&p, (*jwt.HMACSHA)(a.APISecret))
}

func (a *CommonAPI) Shutdown(ctx context.Context) error {
	a.ShutdownChan <- struct{}{}
	return nil
}

func (a *CommonAPI) Session(ctx context.Context) (uuid.UUID, error) {
	return session, nil
}

func (a *CommonAPI) Closing(ctx context.Context) (<-chan struct{}, error) {
	return make(chan struct{}), nil // relies on jsonrpc closing
}
