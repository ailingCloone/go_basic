package svc

import (
	"errors"
	"nrs_customer_module_backend/internal/config"
	"nrs_customer_module_backend/internal/middleware"

	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config                             config.Config
	BasicAuthMiddleware                rest.Middleware
	BeforeLoginTokenCheckingMiddleware rest.Middleware
	AfterLoginTokenCheckingMiddleware  rest.Middleware
	SplashTokenCheckingMiddleware      rest.Middleware
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	svcCtx := &ServiceContext{
		Config: c,
	}
	if err := svcCtx.WithBasicAuthMiddleware(); err != nil {
		return nil, err
	}
	if err := svcCtx.WithBeforeLoginTokenCheckingMiddleware(); err != nil {
		return nil, err
	}
	if err := svcCtx.WithAfterLoginTokenCheckingMiddleware(); err != nil {
		return nil, err
	}
	if err := svcCtx.WithSplashTokenCheckingMiddleware(); err != nil {
		return nil, err
	}
	return svcCtx, nil
}

func (s *ServiceContext) WithBasicAuthMiddleware() error {
	// Initialize your middleware here
	basicAuthMiddleware := middleware.NewBasicAuthMiddleware()
	if basicAuthMiddleware == nil {
		return errors.New("failed to initialize basic auth middleware")
	}
	s.BasicAuthMiddleware = basicAuthMiddleware.Handle
	return nil
}
func (s *ServiceContext) WithBeforeLoginTokenCheckingMiddleware() error {
	// Initialize your middleware here
	beforeLoginTokenCheckingMiddleware := middleware.NewBeforeLoginTokenCheckingMiddleware()
	if beforeLoginTokenCheckingMiddleware == nil {
		return errors.New("failed to initialize before login token checking middleware")
	}

	// Assign the middleware handler to the service context
	s.BeforeLoginTokenCheckingMiddleware = beforeLoginTokenCheckingMiddleware.Handle
	return nil
}

func (s *ServiceContext) WithAfterLoginTokenCheckingMiddleware() error {
	// Initialize your middleware here
	afterLoginTokenCheckingMiddleware := middleware.NewAfterLoginTokenCheckingMiddleware()
	if afterLoginTokenCheckingMiddleware == nil {
		return errors.New("failed to initialize after login token checking middleware")
	}

	// Assign the middleware handler to the service context
	s.AfterLoginTokenCheckingMiddleware = afterLoginTokenCheckingMiddleware.Handle
	return nil
}

func (s *ServiceContext) WithSplashTokenCheckingMiddleware() error {
	// Initialize your middleware here
	splashTokenCheckingMiddleware := middleware.NewSplashTokenCheckingMiddleware()
	if splashTokenCheckingMiddleware == nil {
		return errors.New("failed to initialize after login token checking middleware")
	}

	// Assign the middleware handler to the service context
	s.SplashTokenCheckingMiddleware = splashTokenCheckingMiddleware.Handle
	return nil
}
