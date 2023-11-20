package businesses

import "github.com/go-chi/chi/v5"

func AddBusinessRouter(r *chi.Mux, controller *BusinessesController) {
	r.Route("/businesses", func(r chi.Router) {
		r.With(controller.globalHelpers.SimpleIdentification).With(checkGetBusinessHomeFoodPayload).Get("/", controller.getBusinessHomeFood)
		r.With(controller.globalHelpers.IdentifyUser).With(checkNewBusinessPayload).Post("/", controller.newBusiness)
		r.With(controller.globalHelpers.IdentifyUser).With(checkGetBusinessByIdPayload).Get("/{business-id}", controller.getBusinessById)
		r.With(controller.globalHelpers.IdentifyUser).With(checkNewBusinessMemberPayload).Post("/{business-id}/members", controller.newMember)
		r.With(controller.globalHelpers.IdentifyUser).With(checkCreateFoodPayload).Post("/{business-id}/foods", controller.createFood)
	})
}
