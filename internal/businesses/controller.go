package businesses

import (
	"net/http"

	mw "github.com/Foody-App-Tech/Main-server/internal/global_middlewares"
)

type BusinessesController struct {
	businessesService *BusinessesService
	globalHelpers     *mw.GlobalMiddlewares
}

func NewBusinessesController(businessesService *BusinessesService, globalHelpers *mw.GlobalMiddlewares) *BusinessesController {
	return &BusinessesController{
		businessesService: businessesService,
		globalHelpers:     globalHelpers,
	}
}

func (c *BusinessesController) newBusiness(w http.ResponseWriter, r *http.Request) {

}
