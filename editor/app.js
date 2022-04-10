const MAP_SIZE = 50

const data = {
  map: null,
  mode: null,
  pickerMonster: ["orc", "ghoul", "skeleton", "thing", "wiz", "spectre"],
  pickerWall: [
    "brick_brown_0",
    "brick_brown_2",
    "brick_brown_7",
    "brick_brown-vines_1",
    "brick_brown-vines_3",
    "brick_dark_0",
    "brick_dark_2",
    "brick_dark_4",
    "brick_gray_1",
    "catacombs_0",
    "catacombs_15",
    "catacombs_2",
    "catacombs_4",
    "church_2",
    "cobalt_stone_9",
    "crystal_wall_0",
    "crystal_wall_1",
    "crystal_wall_2",
    "crystal_wall_6",
    "crystal_wall_7",
    "hell_1",
    "hell_7",
    "hell_8",
    "hive_0",
    "hive_2",
    "lab-metal_0",
    "lab-metal_1",
    "lab-metal_3",
    "lab-metal_5",
    "lab-rock_0",
    "lab-rock_1",
    "lab-stone_0",
    "lab-stone_1",
    "lab-stone_5",
    "lair_1_old",
    "marble_wall_11",
    "marble_wall_2",
    "marble_wall_5",
    "marble_wall_9",
    "orc_4",
    "orc_6",
    "orc_7",
    "pebble_red_3_new",
    "relief_0",
    "relief_1",
    "relief_3",
    "sandstone_wall_2",
    "sandstone_wall_5",
    "slime_4",
    "slime_6",
    "snake_1",
    "snake_2",
    "snake_7",
    "snake_9",
    "stone_2_gray0",
    "stone_2_gray1",
    "stone_black_marked_0",
    "stone_black_marked_1",
    "stone_black_marked_7",
    "stone_dark_0",
    "stone2_brown_3_new",
    "tomb_1",
    "undead_brown_0",
    "undead_brown_3",
    "vault_0",
    "vault_2",
    "volcanic_wall_2",
    "volcanic_wall_3",
    "wall_vines_0",
    "wall_vines_2",
    "wall_vines_5",
    "wall_yellow_rock_0",
    "wall_yellow_rock_1",
  ],
  pickerItem: ["potion", "key_green", "key_red", "key_blue", "meat", "apple", "crystal", "column", "barrel"],
  pickerDoor: ["basic", "key_blue", "key_red", "key_green", "switch"],
  pickerDeco: ["torch", "blood_1", "blood_2", "slime", "grate", "switch", "secret", "exit"],
  selectedMonster: 0,
  selectedWall: 0,
  selectedItem: 0,
  selectedDoor: 0,
  selectedDeco: 0,
  cellTip: "",
  fileHandle: null,
  fileName: "",
  loadingSaving: false,
  playerPos: [1, 1],
  floorColour: [1, 1, 1],
  ceilingColour: [1, 1, 1],

  initApp() {
    this.fileHandle = null
    this.fileName = ""
    this.mode = null
    this.map = new Array(MAP_SIZE)
    for (let x = 0; x < MAP_SIZE; x++) {
      this.map[x] = new Array(MAP_SIZE)
      for (let y = 0; y < MAP_SIZE; y++) {
        this.map[x][y] = newEmptyCell(x, y)
      }
    }
    this.playerPos = [1, 1]
    this.map[1][1].t = "p"
    this.map[1][1].v = "0"
    this.ceilingColour = [1, 1, 1]
    this.floorColour = [1, 1, 1]
  },

  cellClick(x, y, evt) {
    // Single click events
    if (evt.type === "click") {
      if (this.mode == "monster") {
        if (this.map[x][y].t == "w" || this.map[x][y].t == "p") return
        this.map[x][y].v = this.pickerMonster[this.selectedMonster]
        this.map[x][y].t = "m"
        return
      }

      if (this.mode == "item") {
        if (this.map[x][y].t == "w" || this.map[x][y].t == "p") return
        this.map[x][y].v = this.pickerItem[this.selectedItem]
        this.map[x][y].t = "i"
        return
      }

      if (this.mode == "door") {
        if (this.map[x][y].t == "w" || this.map[x][y].t == "p") return
        this.map[x][y].v = this.pickerDoor[this.selectedDoor]
        this.map[x][y].t = "d"
        return
      }

      if (this.mode == "extra") {
        // Can only add extras to walls
        if (this.map[x][y].t != "w") return

        // Secret walls are special
        if (this.pickerDeco[this.selectedDeco] == "secret") {
          this.map[x][y].e = ["secret"]
          return
        }

        // Exit walls are special
        if (this.pickerDeco[this.selectedDeco] == "exit") {
          console.log("!!!!!!!!!!!!!!!!exit")
          this.map[x][y].e = ["exit"]
          return
        }

        // Adds a switch
        if (this.pickerDeco[this.selectedDeco] == "switch") {
          const target = prompt("Enter the target cell for this switch:", "x,y")
          if (!target) return
          const targetParts = target.split(",")
          this.map[x][y].e = ["switch", targetParts[0], targetParts[1]]
          return
        }

        this.map[x][y].e = ["deco", this.pickerDeco[this.selectedDeco]]
        return
      }
      if (this.mode == "player") {
        if (this.map[x][y].t == "w") return

        // If player here rotate them
        if (this.map[x][y].t == "p") {
          let facing = parseInt(this.map[x][y].v) + 1
          if (facing > 3) facing = 0
          this.map[x][y].v = "" + facing
          return
        }
        this.map[this.playerPos[0]][this.playerPos[1]].t = null
        this.map[this.playerPos[0]][this.playerPos[1]].v = null
        this.map[x][y].t = "p"
        this.map[x][y].v = "0"
        this.playerPos = [x, y]
        return
      }

      this.cellWall(x, y)
    }

    // Triggered by mousemove to drag & paint walls
    if (evt.buttons === 1) {
      if (this.mode) return
      this.cellWall(x, y)
    }

    this.cellTip = `${x},${y} (${(x + 1) * 32 - 16}, ${(y + 1) * 32 - 16})`
  },

  cellClear(x, y) {
    if (this.map[x][y].t == "p") {
      return
    }

    this.map[x][y] = newEmptyCell(x, y)
  },

  cellWall(x, y) {
    if (this.map[x][y].t == "p") {
      return
    }
    this.map[x][y].t = "w"
    this.map[x][y].v = this.pickerWall[this.selectedWall]
  },

  getCell(x, y) {
    if (!this.map || !this.map[x] || !this.map[x][y]) return null
    return this.map[x][y]
  },

  // Get the main image for a cell
  getImageForCell(x, y) {
    if (!this.map || !this.map[x] || !this.map[x][y]) return "none"
    const cell = this.map[x][y]
    if (!cell || cell === undefined) return "none"
    if (cell.t == "p") return `url(/extra-gfx/player${cell.v}.png)`
    if (cell.t == "i") return `url(/gfx/items/${cell.v}.png)`
    if (cell.t == "m") return `url(/gfx/monsters/${cell.v}.png)`
    if (cell.t == "w") return `url(/gfx/walls/${cell.v}.png)`
    if (cell.t == "d") return `url(/gfx/doors/${cell.v}.png)`
    return "none"
  },

  // Used for decorations and extra stuff on walls  like secrets & exits
  getOverlayForCell(x, y) {
    if (!this.map || !this.map[x] || !this.map[x][y]) return "none"
    const cell = this.map[x][y]
    if (!cell || cell === undefined) return "none"
    if (cell.e.length > 0) {
      if (cell.e[0] == "deco") {
        return `url(/gfx/decoration/${cell.e[1]}.png)`
      }
      if (cell.e[0] == "secret") {
        return `url(/gfx/decoration/secret.png)`
      }
      if (cell.e[0] == "switch") {
        return `url(/gfx/decoration/switch.png)`
      }
      if (cell.e[0] == "exit") {
        return `url(/gfx/decoration/exit.png)`
      }
    }
  },

  setMode(evt) {
    switch (evt.key) {
      case "m":
        this.mode = "monster"
        break
      case "i":
        this.mode = "item"
        break
      case "d":
        this.mode = "door"
        break
      case "x":
        this.mode = "extra"
        break
      case "p":
        this.mode = "player"
        break
    }
  },

  async saveFile() {
    this.loadingSaving = true
    if (!this.fileHandle) {
      try {
        this.fileHandle = await window.showSaveFilePicker(pickerOpts)
      } catch (e) {
        this.loadingSaving = false
        return
      }
    }

    try {
      const writable = await this.fileHandle.createWritable()
      // remove nulls from map

      await writable.write(
        JSON.stringify({
          cells: this.map,
          floorColour: this.floorColour,
          ceilingColour: this.ceilingColour,
        })
      )
      await writable.close()
    } catch (e) {
      console.log(e)
      this.loadingSaving = false
    }
    this.loadingSaving = false
  },

  async openFile() {
    this.loadingSaving = true
    try {
      ;[this.fileHandle] = await window.showOpenFilePicker(pickerOpts)
    } catch (e) {
      this.loadingSaving = false
      return
    }

    if (this.fileHandle.kind === "file") {
      const file = await this.fileHandle.getFile()
      this.fileName = file.name
      const data = await file.text()
      try {
        const rawFile = JSON.parse(data)
        this.map = rawFile.cells
        this.floorColour = rawFile.floorColour
        this.ceilingColour = rawFile.ceilingColour
        for (let x = 0; x < MAP_SIZE; x++) {
          for (let y = 0; y < MAP_SIZE; y++) {
            if (this.map[x][y].t == "p") {
              this.playerPos = [x, y]
            }
          }
        }
      } catch (e) {
        console.error(e)

        alert("Invalid map file, loading aborted")
      }
    }
    this.loadingSaving = false
  },

  setFloorCeiling() {
    //prompt for colors
    const floor = prompt("Floor colour (r,g,b) 0-1", this.floorColour.join(","))
    floor.split(",").forEach((c, i) => {
      this.floorColour[i] = parseFloat(c)
    })
    const ceil = prompt("Ceiling colour (r,g,b) 0-1", this.ceilingColour.join(","))
    ceil.split(",").forEach((c, i) => {
      this.ceilingColour[i] = parseFloat(c)
    })
    // this.floorColour = this.pickerFloor
    // this.ceilingColour = this.pickerCeiling
  },
}

function newEmptyCell(x, y) {
  return {
    x: x,
    y: y,
    t: null,
    v: null,
    e: [],
  }
}

const pickerOpts = {
  types: [
    {
      description: "JSON files",
      accept: {
        "json/*": [".json"],
      },
    },
  ],
  excludeAcceptAllOption: true,
  multiple: false,
}
