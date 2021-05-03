# fRPC

![](https://img.shields.io/badge/stability-experimental-red)

fRPC is a mod that adds RPC capabilities to Factorio. It provides two new buildings that interact with the circuit network:

- _Sensors_ read values from the circuit network and expose them as metrics on the fRPC HTTP server.
- _Receivers_ broadcast values onto the circuit network from inputs provided to routes on the fRPC HTTP server.

```
sudo docker build -t frpc-factorio:0.0.1 .
```

```
sudo docker run \
  -d -it \
  -p 34197:34197/udp \
  -p 27015:27015/tcp \
  -p 8000:8000 \
  --mount type=bind,source="$(pwd)"/data/map-gen-settings.json,target=/factorio/config/map-gen-settings.json \
  --mount type=bind,source="$(pwd)"/data/map-settings.json,target=/factorio/config/map-settings.json \
  --mount type=bind,source="$(pwd)"/data/server-settings.json,target=/factorio/config/server-settings.json \
  -v "$(pwd)"/saves:/factorio/saves \
  -e ENABLE_SERVER_LOAD_LATEST=false \
  -e SAVE_NAME=name-of-your-save \
  --name frpc-factorio \
  frpc-factorio:0.0.1
```

<!--
Use Grafana Cloud

Download agent locally, otherwise you have to forward ports correctly in Docker.
-->

```
docker run \
  -v /tmp/agent:/etc/agent \
  -v $(pwd)/path/to/agent/config.yml:/etc/agent-config/agent.yaml \
  --entrypoint "/bin/agent -config.file=/etc/agent-config/agent.yaml -prometheus.wal-directory=/etc/agent/data" \
  grafana/agent:v0.13.1
```

<!--

TODO:

- Lua mod: Let player enter and specify network name per sensor
  - Don't need to map signal network IDs to meanings
  - Let people purposely name two sensors to be on same "channel", with union semantics
- Separate gauge and counter in game (makes for confusing graphs otherwise)
- Make mod work in singleplayer (see control.lua note)
- Add nicer sprites
-->
