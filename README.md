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
   go run cmd/server/main.go
   ```

1. Run a gomertime client in another window/tab:

   ```shell
   go run cmd/textclient/main.go
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
4. I want to support multiple-client and persistent server scenarios.
   - ...so we'll move towards client/server architecture early
     - ...and use JSON and WebSockets to learn HTTP-oriented go development, even though those probably aren't the right choices for a high-performance game framework.
**5. Don't get lost in the details
   - ...so we'll minimize following of ECS literature.
   - ...so won't be afraid to take shortcuts around hard problems that real games need to solve.
   - ...so I'll minimize testing because the code structure is likely to change repeatedly due to my unfamilarity with the problem space.
   - ...so the UIs will start as text-only.
   - ...so we'll fail fast (let errors happen) and minimize focus on correctness and robustness.

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
  - ~~Text client is able to connect to websocket server, receive server tick, and display~~
  - ~~Log level can be changed at start and at runtime~~
  - ~~Enforce world boundaries during motion, worldXMin, worldXMax...~~
  - ~~Toroidal world option to wrap around world~~
  - ~~Text client is able to show updating server positions~~
  - ~~Move position summary websockets push to array of objects, separate the x/y components while eliminating string key~~
  - ~~Text client is able to move around its viewport independently~~
  - Split server and client fully

- Shortlist
  - Fix client 0,0 not being same as server 0,0
  - Client can connect to custom hostname/port. Client connection over tailscale tunnel confirmed.
  - Server is headless
  - Utility function to count number of neighbors
  - Rewrite server using `time.NewTicker`
  - Add Acceleration component, including a cyclical function to watch an entity cycle back and forth. Ooo, and then circles.
  - Add Health component
  - Find nearby entities
  - Conway game of life rule implementation

- Server/client
  - Text client is able to send commands/updates to server
  - Client registers its viewport, and server updates distinct caches for each client. Client at least receives some subset of every position in the world, as an optimization.
  - Any clients is able to pause the server
  - Clients are able to vote on pausing the server. When all clients have requested the server be paused, the server pauses until requested by any client to restart.

- Display, UI
  - Limit console vertical to a maximum height
  - Camera follows entity
  - Web canvas output
  - Zoom in/out
  - Text display dynamically resizes when window resizes
  - Don't clear text screen every frame, only when needed. Minimize flashing.
  - Dev console/palette
  - Vertical/elevation support in text somehow

- Core, game loop, golang
  - Pass to examine and respond to errors; add error handling throughout
  - Buildings which transform goods
  - systems as goroutines. first: mover managing position
  - Track wall clock time for stats
  - loop over systems (for _, system := range world.Systems() ?)
  - Cache viewport boundaries (xmin/xmax/ymin/ymax)
  - Logging sophistication
  - Genericize into IdType uint64
  - Testing
  - CI
  - Benchmarking
  - WebAssembly
  - Optimization of PositionSummary. Perhaps only calculate all positions within a region
  - Split code into multiple modules
  - Data blob for component, independent of entity-specific data
  - Lua engine embedded to write rules
  - Rationalize the float/integer values used in position summary

- Simulation
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
  - Fields, scalar (or vector) values at each position in the world. Magnetic, photon flux, gravity, ?
  - Engineering with a company. Model dev motivation, daily work
  - Reimplement "puffball" pro-forma ledger automation and modeling infrastructure as entities, and model ticks as the process over time.
