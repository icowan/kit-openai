/**
 * @Time : 2023/3/17 3:01 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package kitopenai

import (
	"context"
	"time"

	"github.com/go-kit/log"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) Models(ctx context.Context) (res ResponseModel, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Models",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Models(ctx)
}

func (s *logging) Model(ctx context.Context, modelName string) (res Model, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Model", "modelName", modelName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Model(ctx, modelName)
}

func (s *logging) Completions(ctx context.Context, req CompletionsRequest) (res ResponseCompletions, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Completions",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Completions(ctx, req)
}

func (s *logging) ChatCompletions(ctx context.Context, req ChatCompletionsRequest) (res ResponseChatCompletions, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "ChatCompletions",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.ChatCompletions(ctx, req)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "logging", "kit-openai")
	return func(next Service) Service {
		return &logging{
			logger:  logger,
			next:    next,
			traceId: traceId,
		}
	}
}
