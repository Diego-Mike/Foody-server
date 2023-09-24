package users

import (
	"github.com/go-chi/chi/v5"
)

func AddUserRouter(r *chi.Mux, controller *UserController) {
	r.Route("/users", func(r chi.Router) {
		r.With(controller.globalHelpers.IdentifyUser).Get("/me", controller.getCurrentUser)
	})
}
