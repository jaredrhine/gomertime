package gomertime

func InitMainWorld(controller *Controller) {
	s := controller.world.store

	pos := s.NewComponent("position")
	vel := s.NewComponent("velocity")

	e1 := s.NewEntity("entity")
	e1.AddComponent(pos, &Position{x: 4, y: 4, z: 0})

	homebase := s.NewEntity("homebase")
	homebase.AddComponent(pos, &Position{x: 10, y: 10, z: 0})

	origin := s.NewEntity("origin")
	origin.AddComponent(pos, &Position{x: 0, y: 0, z: 0})

	mover1 := s.NewEntity("mover")
	mover1.AddComponent(pos, &Position{x: -5, y: -2, z: 0})
	mover1.AddComponent(vel, &Velocity{x: 0.25, y: -0.1, z: 0})

	mover2 := s.NewEntity("mover")
	mover2.AddComponent(pos, &Position{x: 2, y: -2, z: 0})
	mover2.AddComponent(vel, &Velocity{x: 3, y: 0, z: 0})

	mover3 := s.NewEntity("mover")
	mover3.AddComponent(pos, &Position{x: 4, y: -4, z: 0})
	mover3.AddComponent(vel, &Velocity{x: 0, y: -0.75, z: 0})

	whacky1 := s.NewEntity("whacky")
	whacky1.AddComponent(pos, &Position{x: -6, y: -6, z: -1})
}

func InitDevWorld(controller *Controller) {
	s := controller.world.store

	pos := s.NewComponent("position")

	origin := s.NewEntity("origin")
	origin.AddComponent(pos, &Position{x: 0, y: 0, z: 0})

	e1 := s.NewEntity("entity")
	e1.AddComponent(pos, &Position{x: 3, y: 0, z: 0})

	e2 := s.NewEntity("entity")
	e2.AddComponent(pos, &Position{x: 1, y: 1, z: 0})
}

func InitSingleEntityWorld(controller *Controller) {
	s := controller.world.store

	pos := s.NewComponent("position")

	origin := s.NewEntity("origin")
	origin.AddComponent(pos, &Position{x: 0, y: 0, z: 0})
}
