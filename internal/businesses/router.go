package businesses

import "github.com/go-chi/chi/v5"

func AddBusinessRouter(r *chi.Mux, controller *BusinessesController) {
	r.Route("/businesses", func(r chi.Router) {
		r.With(controller.globalHelpers.IdentifyUser).Post("/", controller.newBusiness)
	})
}
