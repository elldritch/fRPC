# fRPC design

Everything in Factorio, including mod behavior, is deterministic. Factorio requires this as a design constraint in order to do parallel, lockstep simulations independently on each client in multiplayer.

Factorio enforces this invariant on mods by executing their Lua scripts in a highly constrained sandbox. In particular, standard library functions that might introduce non-determinism (e.g. file I/O, reading the system clock, listening to network sockets, etc.) are inaccessible.

In order to get around this, we work inside of Factorio's sandbox and run a sidecar alongside a Factorio server in order to provide functionality. The Factorio mod writes outbound data into log files and interprets inbound commands via RCON, while the sidecar hosts the HTTP server and implements handler logic.

## Sensors

Sensors provide an HTTP interface for metrics exported from within Factorio.

### User interface

Within Factorio, sensors are identified by the ID of the circuit network to which they are connected.

Externally, the HTTP API provides these routes:

```
# List basic information about sensors.
GET /sensors

# List detailed information about a sensor's current values.
GET /sensors/:sensor/values
```

### Implementation

Within Factorio, fRPC uses [`game.write_file`](https://lua-api.factorio.com/0.17.79/LuaGameScript.html#LuaGameScript.write_file) to append the circuit network values of each sensor to a series of log files. These writes occur every tick that there is new information (up to 60 times per second).

The readings for every 60 ticks are grouped into a "bucket" and logged in the same file. This file is named `frpc_sensors_$TICK.log`, where `$TICK` is the first tick of the bucket. This is the simplest way to enable log file rotation from within Factorio.

Each line in this file is a JSON value containing a sensor reading:

```typescript
type Line = Reading[];

type Reading = {
  tick: int;
  network_id: int;
  signals: Signal[];
}

type Signal = {
  signal: {
    type: "item" | "fluid" | "virtual";
    name?: string;
  }
  count: int;
}
```

The sidecar reads these log files and updates the current values returned by the HTTP API. It deletes old log files after reading them.
