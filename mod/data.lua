local sensor = table.deepcopy(data.raw["lamp"]["small-lamp"])

sensor.name = "frpc-sensor"
sensor.minable.result = "frpc-sensor"

data:extend({
  -- Entity.
  -- TODO: add custom icon.
  sensor,

  -- Item.
  -- TODO: add custom icon.
  {
    type = "item",
    name = "frpc-sensor",
    icon = "__base__/graphics/icons/constant-combinator.png",
    icon_size = 32,
    place_result = "frpc-sensor",
    subgroup = "circuit-network",
    order = "s",
    stack_size = 10,
  }

  -- TODO: Recipe.


  -- TODO: Technology.
})
