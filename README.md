# Gomertime: Baby's first ECS simulation in go

Author: Jared Rhine <jared@wordzoo.com>

Last update: April 2023

## Usage

1. Install `go` using your preferred approach.

1. Get the codebase from github:

   ```shell
    go version
    git clone https://github.com/jaredrhine/gomertime
    cd gomertime
    go mod download
   ```

1. Run the gomertime server in one window/tab:

   ```shell
   go run cmd/gomertime/main.go
   ```

1. Run a gomertime client in another window/tab:

   ```shell
   go run cmd/gomertime/serverdump/main.go
   ```

Text UI keys:

- `q` - quit
- `escape` - quit
- `control+c` - quit
- `space` - pause
- `d` - dev/debug screen
- arrow keys - move the viewport up/down/left/right

## Design principles

This codebase and line of work is undertaken with these goals in mind:

1. This is a learning exercise, not an attempt to build a real game or community.
   - ...so the whole codebase will be modern go (as written and evolved by a beginner).
2. I want to learn go and in particular concurrency patterns using go channel idioms.
3. I want to learn a bit about ECS (Entity/Component/System) architecture commonly used in games.
   - ...so we'll recreate an ECS architecture from first principles, rather than using an existing ECS library.
   - ...so we'll genericize the implementation to support two or more simulations within this one codebase.
   - ...so we'll not start with an "archetype"-based ECS architecture, which is commonly done as a performance optimization to group entities for faster lookup.
4. Don't get lost in the details
   - Just get an end-to-end simulation running. Don't be afraid to take shortcuts around hard problems that real games need to solve.
   - Focus on text-only interfaces.
5. I want to support multiple-client and persistent server scenarios.
   - ...so we'll move towards client/server architecture early
     - ...and use JSON and WebSockets to learn HTTP-oriented go development, even though those probably aren't the right choices for a high-performance game framework.

## Sketch architecture

A "component" is a typed data bag, such as Position, Velocity, Food, Health that applies to one or more

An "entity" is really just a key/ID to represents an object or container which has multiple components.

A "system" is a code module which performs operations on subsets of components (say, "update position for all entities that have a velocity component")

Our first components will include:

- Position
- Velocity
- Health

Both Position and Velocity are modeled as a three-tuple of float values.

For Position, the 3-tuple represents the positive or negative X, Y, and Z cartesian coordinates of the entity located in a world centered on (0, 0, 0). The units for coordinate values are meters.

For Velocity, the 3-tuple represents a 3D vector, pointed in a specific direction. The units for vector values are meters per second.

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
  - Rationalize the float/integer values used in position summary
  - Utility function to count number of neighbors
  - Implement
  - Rewrite using `time.NewTicker`
  - Add Health component
  - Find nearby entities
  - Conway game of life rule implementation

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
  - Zoom in/out
  - Don't clear text screen every frame, only when needed. Minimize flashing.
  - Vertical/elevation support in text somehow
  - Dev console/palette

- Core, game loop, golang
  - Buildings which transform goods
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
  - Lua engine embedded to write rules

- Simulation
  - Enforce world boundaries during motion, worldXMin, worldXMax...
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
  - Engineering with a company. Model dev motivation, daily work
  - Reimplement "puffball" pro-forma ledger automation and modeling infrastructure as entities, and model ticks as the process over time.
