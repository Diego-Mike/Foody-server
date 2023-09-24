package login

import (
	"github.com/go-chi/chi/v5"
)

func AddLoginRoutes(r *chi.Mux, controller *LoginController) {

	r.Route("/login", func(r chi.Router) {
		r.Get("/oauth/google", controller.googleLogin)
		r.Get("/oauth/facebook", controller.facebookLogin)
		//r.Get("/oauth/tiktok", controller.tiktok)
		//r.Get("/oauth/instagram", controller.instagram)
	})

	r.Post("/access-token", controller.accessToken)
	r.With(controller.globalHelpers.IdentifyUser).Post("/logout", controller.logout)

}
