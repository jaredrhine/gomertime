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

- Done
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
  - ~~Dynamic terminal height/width~~
  - ~~Display representation of world~~
  - ~~Different characters for different entities~~
  - ~~First test for position to text screen mapping~~
  - ~~Display viewport can move around world~~
  - ~~Move WorldStore from controller to world~~
  - ~~Position updates based on velocity each tick~~
  - ~~Split out main package from pkg code~~
  - ~~Logging file is set by default, stderr avoided~~
  - ~~WebSocket server framework, listens on socket~~

- Shortlist
  - Utility function to count number of neighbors

- Server/client
  - Move position summary websockets push to array of objects, separate the x/y components while eliminating string key
  - Text client that is able to connect to websocket server, receive position pushes, and display
  - Text client is able to send commands/updates to server
  - Text client display for multiple simultaneous clients
  - Text client registers a viewport, and updates a viewport
  - Register and calculate multiple viewports, supporting multiple clients looking at different places of the world. Clients register their viewport for the server to precalculate. Server pushes only their viewport
  - Split server and client fully

- Display, UI
  - Limit console vertical to a maximum height
  - Camera follows entity
  - Web canvas output
  - Don't clear text screen every frame, only when needed. Minimize flashing.
  - Vertical/elevation support in text somehow
  - Dev console/palette

- Core, game loop, golang
  - Find nearby entities
  - Buildings which transform goods
  - Rewrite using `time.NewTicker`
  - Log level can be changed at start and at runtime(?)
  - systems as goroutines. first: mover managing position
  - loop over systems (for _, system := range world.Systems() ?)
  - Cache viewport boundaries (xmin/xmax/ymin/ymax)
  - Logging sophistication
  - Genericize into IdType uint64
  - Testing
  - Benchmarking
  - WebAssembly
  - Optimization of PositionSummary. Perhaps only calculate all positions within a region
  - Split code into multiple modules
  - Data blob for component, independent of entity-specific data
  - Turn off worldTickMax by default

- Simulation
  - Enforce world boundaries during motion, worldXMin, worldXMax...
  - Conway game of life rule implementation
  - Food consumption
  - Pathfinding
  - Force and acceleration
  - Friction
  - Destructability
  - Optional rule/flag to not allow movement into same physical space?
  - Toxins
  - Light source
  - Mass of iron
  - Multiple components at same position (works now)
  - Entity size/shape somehow. Entities can spread over multiple coordinates. Center of mass position.
  - Edge of world handling. Bounce off walls. Hard stop.
  - Toroidal world option to wrap around world
  - Fields, scalar (or vector) values at each position in the world. Magnetic, photon flux, gravity, ?
