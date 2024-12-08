package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/object-cache/services/store/pkg/api"
)

type Config struct {
	JWKSURL string
}

func NewWithConfig(ctx context.Context, cfg Config) (echo.MiddlewareFunc, error) {
	jwkc, err := jwk.NewCache(
		ctx,
		httprc.NewClient(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create jwk cache: %w", err)
	}

	err = jwkc.Register(
		ctx,
		cfg.JWKSURL,
		jwk.WithMaxInterval(24*time.Hour),
		jwk.WithMinInterval(15*time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register jwk cache: %w", err)
	}

	log.Ctx(ctx).Info().Str("url", cfg.JWKSURL).Msg("jwk cache registered")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			cached, err := jwkc.CachedSet(cfg.JWKSURL)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("failed to parse token")
				return c.JSON(http.StatusUnauthorized, api.Error{
					Message: "invalid token",
				})
			}

			tok, err := jwt.ParseRequest(c.Request(), jwt.WithKeySet(cached), jwt.WithHeaderKey("Authorization"))
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("failed to parse token")
				return c.JSON(http.StatusUnauthorized, api.Error{
					Message: "invalid token",
				})
			}

			log.Ctx(ctx).Info().Any("tok", tok).Msg("token")

			return next(c)
		}
	}, nil
}
