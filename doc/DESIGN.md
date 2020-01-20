# fRPC design

Everything in Factorio, including mod behavior, is deterministic. Factorio requires this as a design constraint in order to do parallel, lockstep simulations independently on each client in multiplayer.

Factorio enforces this invariant on mods by executing their Lua scripts in a highly constrained sandbox. In particular, standard library functions that might introduce non-determinism (e.g. file I/O, reading the system clock, listening to network sockets, etc.) are inaccessible.

In order to get around this, we work inside of Factorio's sandbox and run a sidecar alongside a Factorio server in order to provide functionality. The Factorio mod writes outbound data into log files and interprets inbound commands via RCON, while the sidecar hosts the HTTP server and implements handler logic.

## Sensors

Sensors provide an HTTP interface for metrics exported from within Factorio.

### User interface

Within Factorio, sensors are assigned a UUID when they're created. Players can also assign a human-friendly name string to the sensor. Names must be unique.

Externally, the HTTP API provides these routes:

```
# List basic information about sensors, including a map of sensor names to IDs.
GET /sensors

# List detailed information about a sensor's current values.
GET /sensors/:sensor/values
```

### Implementation

Within Factorio, fRPC uses [`game.write_file`](https://lua-api.factorio.com/0.17.79/LuaGameScript.html#LuaGameScript.write_file) to write the circuit network values of each sensor to a JSON log file once per second. These files are named `fRPC_log_$TIMESTAMP.json`, where `$TIMESTAMP` is the timestamp at which the sensor values were sampled truncated to second resolution.

The format of these files is:

```typescript
type UUID = string;

type File = {
  names: {
    [name: string]: UUID
  }
  values: {
    [sensor: UUID]: {
      [channel: string]: number;
    }
  }
};
```

The sidecar reads these log files every second and updates the current values returned by the HTTP API. It deletes log files after reading them.
