package router

import "github.com/labstack/echo/v4"

type Route struct {
	Method      string
	Path        string
	Handler     echo.HandlerFunc
	Middlewares []echo.MiddlewareFunc
}

type RouteRegistrar interface {
	RegisterRoutes() []Route
}

func Register(e *echo.Echo, registrars ...RouteRegistrar) {
	for _, r := range registrars {
		for _, route := range r.RegisterRoutes() {
			switch route.Method {
			case "GET":
				e.GET(route.Path, route.Handler, route.Middlewares...)
			case "POST":
				e.POST(route.Path, route.Handler, route.Middlewares...)
			case "PUT":
				e.PUT(route.Path, route.Handler, route.Middlewares...)
			case "DELETE":
				e.DELETE(route.Path, route.Handler, route.Middlewares...)
			case "PATCH":
				e.PATCH(route.Path, route.Handler, route.Middlewares...)
			}
		}
	}
}
