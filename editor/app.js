const MAP_SIZE = 50;

const data = {
  map: null,
  mode: null,

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
      if (this.mode == "m") {
        this.map[x][y].monster = "ghoul";
        return;
      }
      this.cellWall(x, y, 3);
    }

    if (evt.buttons === 1) {
      this.cellWall(x, y, 3);
    }
  },

  cellClear(x, y) {
    this.map[x][y] = newEmptyCell(x, y);
  },

  cellWall(x, y, i) {
    this.map[x][y].monster = null;
    this.map[x][y].item = null;
    this.map[x][y].extra = null;
    this.map[x][y].wall = i;

    //if (x % 3 === 0) {
    this.map[x][y].extra = "sdsssdssdd";
    //}
  },

  getCell(x, y) {
    return this.map[x][y];
  },

  getImageForCell(x, y) {
    const cell = this.map[x][y];
    if (cell.wall) {
      return "url(/gfx/walls/1.png)";
    } else if (cell.monster) {
      return "url(/gfx/monsters/ghoul.png)";
    } else if (cell.item) {
      return "item";
    }
    return "none";
  },

  getOverlayForCell(x, y) {
    const cell = this.map[x][y];
    if (cell.extra) {
      return "url(/img/sprites/ball.png)";
    }
  },
};

function newEmptyCell(x, y) {
  return {
    x: x,
    y: y,
    wall: null,
    monster: null,
    item: null,
    extra: null,
  };
}
