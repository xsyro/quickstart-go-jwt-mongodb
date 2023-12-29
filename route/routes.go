package route

import (
	"context"
	"quickstart-go-jwt-mongodb/controllers"
	"quickstart-go-jwt-mongodb/internal"
	"quickstart-go-jwt-mongodb/server"
)

func Routes(httpHandler server.RequestHandler, database internal.MongoDatabase, ctx context.Context) {

	httpHandler.ControllerRegistry(controllers.Homepage(ctx))
	httpHandler.ControllerRegistry(controllers.HealthCheck(ctx))

	httpHandler.ControllerRegistry(controllers.Authenticate(database, ctx))
	httpHandler.ControllerRegistry(controllers.RefreshToken(database, ctx))
	httpHandler.ControllerRegistry(controllers.CreateAccount(database, ctx))

	httpHandler.ControllerRegistry(controllers.SecuredRole1Only(database, ctx))
	httpHandler.ControllerRegistry(controllers.SecuredRole2Only(database, ctx))
	httpHandler.ControllerRegistry(controllers.SecuredRole1And2Only(database, ctx))
}
