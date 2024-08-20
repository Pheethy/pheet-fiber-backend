package handler

import (
	"context"
	"fmt"
	"net/http"
	"pheet-fiber-backend/auth"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/middleware"
	"pheet-fiber-backend/service/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
)

type middlewareHandler struct {
	cfg      config.Iconfig
	middleUs middleware.IMiddlewareUsecase
}

func NewMiddlewareHandler(cfg config.Iconfig, middleUs middleware.IMiddlewareUsecase) middleware.ImiddlewareHandler {
	return middlewareHandler{
		cfg:      cfg,
		middleUs: middleUs,
	}
}

func (m middlewareHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, PUT, PATCH, HEAD, DELETE",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (m middlewareHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "👽 ${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Bangkok/Asia",
	})
}

func (m middlewareHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		mapClaims, err := auth.ParseToken(m.cfg.Jwt(), token)
		if err != nil {
			return fiber.NewError(http.StatusUnauthorized, err.Error())
		}

		claims := mapClaims.Claims

		if !m.middleUs.FindAccessToken(ctx, claims.Id, token) {
			return fiber.NewError(http.StatusUnauthorized, "no permission to access")
		}

		c.Locals("user_id", claims.Id)
		c.Locals("role_id", claims.RoleId)
		return c.Next()
	}
}

func (m middlewareHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("user_id")
		if c.Params("user_id") != userId {
			return fiber.NewError(http.StatusBadRequest, "never gonna give you up")
		}
		return c.Next()
	}
}

func (m middlewareHandler) Authorize(expectedRoleId ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		userRoleId, ok := c.Locals("role_id").(int)
		if !ok {
			return fiber.NewError(http.StatusUnprocessableEntity, "cast role_id to int failed.")
		}

		roles, err := m.middleUs.FetchRoles(ctx)
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}

		sum := 0
		for _, val := range expectedRoleId {
			sum += val
		}

		expectedValBinary := utils.ConvertBinary(sum, len(roles))
		userValBinary := utils.ConvertBinary(userRoleId, len(roles))

		for index := range userValBinary {
			if userValBinary[index]&expectedValBinary[index] == 1 {
				return c.Next()
			}
		}

		return fiber.NewError(http.StatusUnauthorized)
	}
}

func (m middlewareHandler) ApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-API-KEY")
		if _, err := auth.ParseAPIKey(m.cfg.Jwt(), key); err != nil {
			return fiber.NewError(http.StatusInternalServerError, "API-KEY is invalid.")
		}
		return c.Next()
	}
}

func (m middlewareHandler) SetTracer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var span opentracing.Span
		ctx := c.UserContext()
		spanName := fmt.Sprintf("%s %s %s", string(c.Context().Request.URI().Scheme()), c.Method(), c.Path())
		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.GetReqHeaders()))
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}
		switch err {
		case nil:
			span = opentracing.StartSpan(spanName, ext.RPCServerOption(spanCtx))
		case opentracing.ErrSpanContextNotFound:
			span, ctx = opentracing.StartSpanFromContext(ctx, spanName)
		default:
			logrus.Println("error default")
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}
		defer span.Finish()

		c.SetUserContext(ctx)
		// Proceed to the next handler
		err = c.Next()

		m.setTagByFiber(span, c)
		m.setLogByFiber(span, c)

		if err != nil {
			m.setError(span, c, err)
		} else {
			span.SetTag("error", false)
			span.SetTag("http.status_code", c.Response().StatusCode())
		}

		return nil
	}
}

// Note: Fiber doesn't support parameter names directly, so this function is omitted for brevity

func (m middlewareHandler) setTagByFiber(span opentracing.Span, c *fiber.Ctx) {
	span.SetTag("host", c.Hostname())
	span.SetTag("User-Agent", c.Get("User-Agent"))
	span.SetTag("http.method", c.Method())
	span.SetTag("http.url", c.OriginalURL())
}

func (m middlewareHandler) setLogByFiber(span opentracing.Span, c *fiber.Ctx) {
	span.LogFields(
		log.String("querystring", c.Context().QueryArgs().String()),
	)
}

func (m middlewareHandler) setError(span opentracing.Span, c *fiber.Ctx, err error) {
	isError := err != nil && c.Response().StatusCode() >= http.StatusBadRequest
	span.SetTag("error", isError)
	if isError {
		span.SetTag("http.status_code", c.Response().StatusCode())
		span.LogFields(log.Message(err.Error()))
	}
}
