const MAP_SIZE = 50

const data = {
  map: null,
  mode: null,
  pickerMonster: ["ghoul", "skeleton", "thing"],
  pickerWall: wallIndex,
  pickerItem: ["potion"],
  pickerDoor: ["basic", "blue_key", "red_key", "green_key", "switch"],
  pickerDeco: ["torch", "blood_1", "blood_2", "slime", "grate", "switch", "secret"],
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
  },

  cellClick(x, y, evt) {
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
        if (this.map[x][y].t != "w") return
        if (this.pickerDeco[this.selectedDeco] == "secret") {
          this.map[x][y].e = ["secret"]
          return
        }
        if (this.pickerDeco[this.selectedDeco] == "switch") {
          const target = prompt("Enter the switch's target cell:", "x,y")
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

    if (evt.buttons === 1) {
      if (this.mode) return
      this.cellWall(x, y)
    }

    this.cellTip = `${x},${y}`
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

  getImageForCell(x, y) {
    if (!this.map || !this.map[x] || !this.map[x][y]) return "none"
    const cell = this.map[x][y]
    if (!cell || cell === undefined) return "none"
    if (cell.t == "p") return `url(/gfx/player${cell.v}.png)`
    if (cell.t == "i") return `url(/gfx/items/${cell.v}.png)`
    if (cell.t == "m") return `url(/gfx/monsters/${cell.v}.png)`
    if (cell.t == "w") return `url(/gfx/walls/${cell.v}.png)`
    if (cell.t == "d") return `url(/gfx/doors/${cell.v}.png)`
    return "none"
  },

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

      await writable.write(JSON.stringify(this.map))
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
        this.map = JSON.parse(data)
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
