package application


type App struct {
	Server
	Log  logger.Interface
}

func New(logger logger.Interface) *App {
	return &App{Log: logger}
}
