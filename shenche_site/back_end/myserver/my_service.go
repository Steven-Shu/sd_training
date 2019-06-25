package myserver

import "github.com/drone/routes"

type MyService interface {
	RegisterServices(m *routes.RouteMux)
}
