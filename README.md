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

## Coordinate transform

```text

```

## Work plan

- ~~Project name~~
- ~~Basic go module~~
- ~~World tick update loop~~
- ~~Terminal clear and display shell~~
- ~~Prototype core with keyboard, display, world tick integrated~~
- ~~Detect keyboard press async~~
- ~~'q' to exit~~
- ~~Pause key~~
- ~~Decide that components will belong to single entities, rather than being able to be shared between entities~~
- ~~id generator (sequence for entity probably)~~
- ~~Component~~
- ~~Position, velocity components stored as 3D float~~
- ~~Entity~~
- ~~Temporary hard-coded world/entity creation/setup~~
- ~~Fake/inefficient ECS storage backend; defer archetype implementation~~
- ~~Debug dump state~~
- ~~Lookup console height/width on start~~
- ~~Position lookup/query~~
- ~~Display representation of world~~
- ~~First test for position to text screen mapping~~
- ~~Move WorldStore from controller to world~~
- Split out main package from pkg code
- Position updates based on velocity each tick
- WebSocket server framework, listens on socket
- Cache viewport boundaries (xmin/xmax/ymin/ymax)
- loop over systems (for _, system := range world.Systems() ?)
- Display viewport can move around world
- Different characters for different entities
- Limit console vertical height
- Dynamic terminal height/width
- Data blob for component, independent of entity-specific data
- Don't clear text screen every frame, only when needed. Minimize flashing.
- systems as goroutines. first: mover managing position
- Conway game of life rule implementation
- Dev console/palette
- Acceleration
- Genericize into IdType uint64
- Buildings which transform goods
- Rewrite using `time.NewTicker`
- Edge of world handling. Bounce off walls. Hard stop.
- Multiple components at same position
- Find nearby items
- Benchmarking
- Pathfinding
- Force and Acceleration
- Split into multiple modules
- Food consumption
- Mass of iron
- Light source
- Logging file is set by default, stderr avoided, tuned
- Camera follows entity
- Split server and client
- Register and calculate multiple viewports, supporting multiple clients looking at different places of the world. Clients register their viewport for the server to precalculate.
- WebAssembly
- Canvas output
- Optimization of PositionSummary. Perhaps only calculate all positions within a region
- Entity size/shape somehow. Entities can spread over multiple coordinates. Center of mass position.
- Friction
- Destructability
- Toxins
- Vertical support in text somehow
