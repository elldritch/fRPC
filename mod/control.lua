local function deduplicate(list)
  local unique_values = {}
  for _, v in ipairs(list) do
    unique_values[v] = true
  end

  local deduplicated = {}
  for k, _ in pairs(unique_values) do
    table.insert(deduplicated, k)
  end

  return deduplicated
end

-- NOTE: these connector IDs are actually not unique -- all are 1 except
-- combinator_output, which is 2.
local connector_ids = {
  defines.circuit_connector_id.accumulator,
  defines.circuit_connector_id.constant_combinator,
  defines.circuit_connector_id.container,
  defines.circuit_connector_id.programmable_speaker,
  defines.circuit_connector_id.rail_signal,
  defines.circuit_connector_id.rail_chain_signal,
  defines.circuit_connector_id.roboport,
  defines.circuit_connector_id.storage_tank,
  defines.circuit_connector_id.wall,
  defines.circuit_connector_id.electric_pole,
  defines.circuit_connector_id.inserter,
  defines.circuit_connector_id.lamp,
  defines.circuit_connector_id.combinator_input,
  defines.circuit_connector_id.combinator_output,
  defines.circuit_connector_id.offshore_pump,
  defines.circuit_connector_id.pump
}

local unique_connector_ids = deduplicate(connector_ids)

local wire_types = {
  defines.wire_type.red,
  defines.wire_type.green
}

script.on_nth_tick(1, function (e)
  local sample = {
    tick = e.tick,
    values = {}
  }
  for _, surface in pairs(game.surfaces) do
    local sensors = surface.find_entities_filtered({name = "frpc-sensor"})
    for _, sensor in ipairs(sensors) do
      for _, wire_type in ipairs(wire_types) do
        for _, connector_id in ipairs(unique_connector_ids) do
          local network = sensor.get_circuit_network(wire_type, connector_id)
          if network ~= nil and network.signals ~= nil then
            table.insert(sample.values, {
              network_id = network.network_id,
              signals = network.signals
            })
          end
        end
      end
    end
  end
  if next(sample.values) ~= nil then
    game.write_file("frpc_sensors_".. math.floor(e.tick / 60) ..".log", game.table_to_json(sample) .. "\n", true)
  end
end)
