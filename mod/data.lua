-- TODO: Need another base - lamps turn themselves off now with circuit input.
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
    icon = "__base__/graphics/icons/small-lamp.png",
    icon_size = 32,
    place_result = "frpc-sensor",
    subgroup = "circuit-network",
    order = "s",
    stack_size = 50,
  },

  -- Recipe.
  {
    type = "recipe",
    name = "frpc-sensor",
    ingredients =
    {
      {"copper-cable", 5},
      {"electronic-circuit", 5},
    },
    result = "frpc-sensor",
    enabled = false
  },

  -- Technology.
  -- TODO: add custom icon.
  {
    type = "technology",
    name = "remote-control",
    icon_size = 128,
    icon = "__base__/graphics/technology/circuit-network.png",
    effects =
    {
      {
        type = "unlock-recipe",
        recipe = "frpc-sensor"
      }
    },
    prerequisites = {"circuit-network", "automation-2"},
    unit =
    {
      count = 150,
      ingredients =
      {
        {"automation-science-pack", 1},
        {"logistic-science-pack", 1}
      },
      time = 15
    },
    order = "r-s-f"
  },
})
