package controller

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"controller",
	fx.Provide(
		fx.Annotate(
			NewBattleController,
			fx.As(new(Controller)),
			fx.ResultTags(`group:"controllers"`),
		),
		fx.Annotate(
			NewCharacterController,
			fx.As(new(Controller)),
			fx.ResultTags(`group:"controllers"`),
		),
		fx.Annotate(
			NewSkillController,
			fx.As(new(Controller)),
			fx.ResultTags(`group:"controllers"`),
		),
	),
)

type Controller interface {
	Mount(fiber.Router)
}
