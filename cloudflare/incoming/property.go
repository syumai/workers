package incoming

import (
	"context"

	"github.com/syumai/workers/internal/cfcontext"
)

type Properties struct {
	AsOrganization string
}

func NewProperties(ctx context.Context) *Properties {
	obj := cfcontext.MustExtractIncomingProperty(ctx)
	return &Properties{
		AsOrganization: obj.Get("asOrganization").String(),
	}
}
