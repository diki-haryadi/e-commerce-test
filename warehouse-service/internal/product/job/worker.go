package productJob

import (
	"context"

	"github.com/diki-haryadi/ztools/wrapper"
)

func (j *job) logProductWorker() wrapper.HandlerFunc {
	return func(ctx context.Context, args ...interface{}) (interface{}, error) {
		j.logger.Info("product log job")
		return nil, nil
	}
}
