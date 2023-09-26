package main

import (
	"context"

	"github.com/alecthomas/kong"
	"github.com/farseeingnorthwest/battleground.go/controller"
	"github.com/farseeingnorthwest/battleground.go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

func main() {
	var cli struct {
		Debug  bool
		DSN    string `env:"DATABASE_URL" required:""`
		Addr   string `default:":3000"`
		Static string
	}
	kong.Parse(&cli)

	fx.New(
		storage.Module,
		controller.Module,
		fx.Supply(
			fx.Annotate(
				cli.Addr,
				fx.ResultTags(`name:"addr"`),
			),
			fx.Annotate(
				cli.Debug,
				fx.ResultTags(`name:"debug"`),
			),
			fx.Annotate(
				cli.Static,
				fx.ResultTags(`name:"static"`),
			),
		),
		fx.Provide(
			func(r *storage.CharacterRepository) controller.CharacterRepository {
				return r
			},
			func(r *storage.SkillRepository) controller.SkillRepository {
				return r
			},
			func() *sqlx.DB {
				db, err := sqlx.Connect("postgres", cli.DSN)
				if err != nil {
					panic(err)
				}
				return db
			},
			NewFiberApp,
		),
		fx.Invoke(func(app *fiber.App) {}),
	).Run()
}

func NewFiberApp(params FiberAppParams, lc fx.Lifecycle) *fiber.App {
	app := fiber.New()
	if params.Debug {
		app.Use(logger.New())
	}
	api := app.Group("/api")
	for _, c := range params.Controllers {
		c.Mount(api)
	}
	if params.Static != "" {
		app.Static("/*", params.Static)
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := app.Listen(params.Addr); err != nil {
					log.Error(err)
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			return app.Shutdown()
		},
	})

	return app
}

type FiberAppParams struct {
	fx.In

	Controllers []controller.Controller `group:"controllers"`
	Addr        string                  `name:"addr"`
	Debug       bool                    `name:"debug"`
	Static      string                  `name:"static"`
}
