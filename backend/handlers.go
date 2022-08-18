package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	EnvKeyPrimaryHost = "DETA_SPACE_APP_HOSTNAME"
)

// Config xx
type Config struct {
	Address string
	CodeStoreName string
	CodeStoreKey *string
}

// Server xx
type Server struct {
	conf *Config
	store *CodeStore
	e *echo.Echo
}

// NewServer xx
func NewServer(c *Config) (*Server, error) {
	store, err := NewCodeStore(c.CodeStoreName, c.CodeStoreKey)
	if err != nil {
		return nil, fmt.Errorf("new code store: %w", err)
	}
	e := echo.New()
	e.Use(middleware.Recover())
	s := &Server{
		conf: c,
		store: store,
		e: e,
	}
	s.registerHandlers()
	s.configureCORS()
	return s, nil
}

// Start xx
func (s *Server) Start() error{
	return s.e.Start(s.conf.Address)
}

func (s *Server) registerHandlers() {
	s.e.GET("/", s.HealthCheck)
	s.e.GET("/codes", s.List)
	s.e.POST("/codes", s.Post)
	s.e.GET("/codes/content", s.Get)
	s.e.DELETE("/codes/content", s.Delete)
}

func (s *Server) configureCORS() {
	s.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", fmt.Sprintf("https://%s", os.Getenv(EnvKeyPrimaryHost))},
		AllowCredentials: true,
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
}

func (s *Server) NewInternalServerError() *echo.HTTPError {
	return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
}

// HealthCheck xx
// GET /
func (s *Server) HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

// List xx
// Get /codes
func (s *Server) List(c echo.Context) error {
	codes, err := s.store.List()
	if err != nil {
		s.e.Logger.Printf("failed to list codes: %v\n", err)
		return s.NewInternalServerError()
	}
	return c.JSON(http.StatusOK, codes)
}

// Get xx
// Get code
func (s *Server) Get(c echo.Context) error {
	key := c.QueryParam("key")
	if key == ""{
		return echo.NewHTTPError(http.StatusBadRequest, "no key")
	}
	code, err := s.store.Get(key)
	if err != nil {
		s.e.Logger.Printf("failed to get code with key %s: %v\n", key, err)
		return s.NewInternalServerError()
	}
	if code == nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}
	return c.String(http.StatusOK, code.Content)
}

// Delete xx
func (s *Server) Delete(c echo.Context) error {
	key := c.QueryParam("key")
	if key == ""{
		return echo.NewHTTPError(http.StatusBadRequest, "no key")
	}
	if err := s.store.Delete(key); err != nil{
		s.e.Logger.Printf("failed to delete code with key %s: %v\n", key, err)
		return s.NewInternalServerError()
	}
	return c.String(http.StatusOK, key)
}

// Post xx
func (s *Server) Post(c echo.Context) error {
	key := c.QueryParam("key")
	if key == ""{
		return echo.NewHTTPError(http.StatusBadRequest, "no key")
	}
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		s.e.Logger.Printf("failed to read request body for key %s: %v\n", key, err)
		return s.NewInternalServerError()
	}
	defer c.Request().Body.Close()

	code := &Code{
		Key: key,
		Content: string(body),
	}
	if err := s.store.Put(code); err != nil {
		s.e.Logger.Printf("failed to put code for key %s: %v\n", key, err)
		return s.NewInternalServerError()
	}
	return c.String(http.StatusOK, key)
}
