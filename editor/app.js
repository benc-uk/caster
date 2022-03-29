const MAP_SIZE = 30;

const data = {
  map: null,
  mode: null,
  pickerMonster: ["ghoul", "skeleton", "thing"],
  pickerWall: ["1", "2", "3", "4", "5", "6", "7", "8"],
  pickerItem: ["potion"],
  pickerUsable: ["basic", "blue"],
  pickerDeco: ["torch_1"],
  selectedMonster: 0,
  selectedWall: 0,
  selectedItem: 0,
  selectedUsable: 0,
  selectedDeco: 0,
  tooltip: "",

  initApp() {
    this.map = new Array(MAP_SIZE);
    for (let x = 0; x < MAP_SIZE; x++) {
      this.map[x] = new Array(MAP_SIZE);
      for (let y = 0; y < MAP_SIZE; y++) {
        this.map[x][y] = newEmptyCell(x, y);
      }
    }
  },

  cellClick(x, y, evt) {
    if (evt.type === "click") {
      if (this.mode == "monster") {
        if (this.map[x][y].w) return;
        this.map[x][y].m = this.pickerMonster[this.selectedMonster];
        return;
      }
      if (this.mode == "item") {
        if (this.map[x][y].w) return;
        this.map[x][y].i = this.pickerItem[this.selectedItem];
        console.log(this.map[x][y]);
        return;
      }
      if (this.mode == "door") {
        if (this.map[x][y].w) return;
        this.map[x][y].i = this.pickerItem[this.selectedItem];
        console.log(this.map[x][y]);
        return;
      }
      this.cellWall(x, y);
    }

    if (evt.buttons === 1) {
      if (this.mode) return;
      this.cellWall(x, y);
    }
    this.tooltip = `${x},${y}`;
  },

  cellClear(x, y) {
    this.map[x][y] = newEmptyCell(x, y);
  },

  cellWall(x, y) {
    wall = this.pickerWall[this.selectedWall];
    this.map[x][y].m = null;
    this.map[x][y].i = null;
    this.map[x][y].d = null;
    this.map[x][y].w = wall;
  },

  getCell(x, y) {
    return this.map[x][y];
  },

  getImageForCell(x, y) {
    const cell = this.map[x][y];
    if (cell.w) {
      return `url(/gfx/walls/${cell.w}.png)`;
    } else if (cell.m) {
      return `url(/gfx/monsters/${cell.m}.png)`;
    } else if (cell.i) {
      return `url(/gfx/items/${cell.i}.png)`;
    }
    return "none";
  },

  getOverlayForCell(x, y) {
    const cell = this.map[x][y];
    if (cell.extra) {
      return "url(/gfx/items/ball.png)";
    }
  },

  setMode(evt) {
    switch (evt.key) {
      case "m":
        this.mode = "monster";
        break;
      case "i":
        this.mode = "item";
        break;
      case "u":
        this.mode = "usable";
        break;
      case "d":
        this.mode = "decoration";
        break;
    }
  },

  save() {
    console.log(JSON.stringify(this.map));
  },
};

function newEmptyCell(x, y) {
  return {
    x: x,
    y: y,
    w: null,
    m: null,
    i: null,
    d: null,
  };
}
