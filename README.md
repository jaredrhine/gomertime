# Gomertime: Baby's first ECS simulation in go

Author: Jared Rhine <jared@wordzoo.com>
Last update: April 2023

## Design principles

For this codebase, we're going to try coding an ECS simulation architecture from scratch as a learning exercise (for both ECS and go).

Don't get lost in the details:

- Get a full end-to-end simulation running.
- Stick with a text-only interface initially.

## Sketch architecture

Jared wants to learn:

- go core language skills...
  - ...so the whole codebase will be modern go (as written by a beginner).
- ECS (Entity/Component/System) data-focused simulation/gaming architecture...
  - ...so we'll recreate an ECS architecture from first principles, rather than using an existing ECS library.
  - ...we'll genericize the implementation to support two or more simulations within this one codebase.

We will not start with an "archetype"-based ECS architecture, generally done as a performance optimization to group entities for faster lookup.

ECS provides multiple typed "data bags". This will be implemented as go named structs.

## Core code structure

## Notes

- Very first model
  - Center-of-mass position in 3-d space (but hard-code z position to zero)
  - Velocity
  - Mass of iron ore carried

- Potential components
  - Position
  - Velocity
  - Charge
  - Center of mass position
  - Light
  - Food

- Simulation ideas
  - Engineering with a company. Model dev motivation, daily work
  - Reimplement "puffball" pro-forma ledger automation and modeling infrastructure as entities, and model ticks as the process over time.

- Lua engine embedded to write rules

## Work plan

- ~~Project name~~
- ~~Basic go module~~
- ~~World tick update loop~~
- ~~Terminal clear and display shell~~
- ~~Prototype core with keyboard, display, world tick integrated~~
- ~~Detect keyboard press async~~
- ~~'q' to exit~~
- ~~Pause key~~
- Temporary hard-coded world/entity creation/setup
- Debug dump state
- Fake/inefficient ECS storage backend; defer archetype implementation
- Position lookup/query
- Display representation of world
- Display viewport can move around world
- id generator (sequence for entity probably)
- NewEntity
- loop over systems (for _, system := range world.Systems() ?)
- Center of mass position
- Velocity
- Rewrite using `time.NewTicker`
- Dev console/palette
- Bounce off walls
- Multiple components at same position
- Find nearby items
- Keyboard
- Force -> Acceleration
- Friction
- Go modules
- Mass of iron
- WebAssembly
- Canvas output
- Go channels
- Camera follows entity
