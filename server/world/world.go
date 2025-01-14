package world

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"

	"github.com/aimjel/minecraft/nbt"
	"github.com/dynamitemc/dynamite/server/world/anvil"
)

type World struct {
	nbt      worldData
	Gamemode byte

	overworld *Dimension
	nether    *Dimension
	theEnd    *Dimension
}

func OpenWorld(name string, flat bool) (*World, error) {
	f, err := os.Open(name + "/level.dat")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var wrld World
	if err = loadWorldData(f, &wrld.nbt); err != nil {
		return nil, fmt.Errorf("%v loading world level data", err)
	}

	wrld.overworld = NewDimension("minecraft:overworld", anvil.NewReader(name+"/region/", name+"/entities/"))
	wrld.nether = NewDimension("minecraft:the_nether", anvil.NewReader(name+"/DIM-1/region/", name+"/DIM-1/entities/"))
	wrld.theEnd = NewDimension("minecraft:the_end", anvil.NewReader(name+"/DIM1/region/", name+"/DIM1/entities/"))
	if flat {
		wrld.overworld.generator = &FlatGenerator{}
	}

	return &wrld, nil
}

func (w *World) Seed() int64 {
	return w.nbt.Data.WorldGenSettings.Seed
}

func (w *World) Spawn() (x, y, z int32, angle float32) {
	return w.nbt.Data.SpawnX, w.nbt.Data.SpawnY, w.nbt.Data.SpawnZ, w.nbt.Data.SpawnAngle
}

func (w *World) Overworld() *Dimension {
	return w.overworld
}

func (w *World) Nether() *Dimension {
	return w.nether
}

func (w *World) TheEnd() *Dimension {
	return w.theEnd
}

func (w *World) LoadSpawnChunks(rd int32) (success int) {
	ow := w.Overworld()
	s := 0
	for x := -rd; x < rd; x++ {
		for z := -rd; z < rd; z++ {
			if _, err := ow.Chunk(x, z); err == nil {
				s++
			}
		}
	}
	return s
}

func (w *World) Gamerules() map[string]string {
	return w.nbt.Data.GameRules
}

func loadWorldData(f *os.File, wNbt *worldData) error {
	gzipRd, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(gzipRd); err != nil {
		return err
	}

	return nbt.Unmarshal(buf.Bytes(), wNbt)
}
