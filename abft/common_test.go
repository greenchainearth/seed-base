package abft

import (
	"github.com/greenchainearth/seed-base/inter/idx"
	"github.com/greenchainearth/seed-base/inter/pos"
	"github.com/greenchainearth/seed-base/kvdb"
	"github.com/greenchainearth/seed-base/kvdb/memorydb"
	"github.com/greenchainearth/seed-base/seed"
	"github.com/greenchainearth/seed-base/utils/adapters"
	"github.com/greenchainearth/seed-base/vecfc"
)

type applyBlockFn func(block *seed.Block) *pos.Validators

// TestSeed extends Seed for tests.
type TestSeed struct {
	*IndexedSeed

	blocks map[idx.Block]*seed.Block

	applyBlock applyBlockFn
}

// FakeSeed creates empty abft with mem store and equal weights of nodes in genesis.
func FakeSeed(nodes []idx.ValidatorID, weights []pos.Weight, mods ...memorydb.Mod) (*TestSeed, *Store, *EventStore) {
	validators := make(pos.ValidatorsBuilder, len(nodes))
	for i, v := range nodes {
		if weights == nil {
			validators[v] = 1
		} else {
			validators[v] = weights[i]
		}
	}

	openEDB := func(epoch idx.Epoch) kvdb.DropableStore {
		return memorydb.New()
	}
	crit := func(err error) {
		panic(err)
	}
	store := NewStore(memorydb.New(), openEDB, crit, LiteStoreConfig())

	err := store.ApplyGenesis(&Genesis{
		Validators: validators.Build(),
		Epoch:      FirstEpoch,
	})
	if err != nil {
		panic(err)
	}

	input := NewEventStore()

	config := LiteConfig()
	lch := NewIndexedSeed(store, input, &adapters.VectorToDagIndexer{vecfc.NewIndex(crit, vecfc.LiteConfig())}, crit, config)

	extended := &TestSeed{
		IndexedSeed: lch,
		blocks:      map[idx.Block]*seed.Block{},
	}

	blockIdx := idx.Block(0)

	err = extended.Bootstrap(seed.ConsensusCallbacks{
		BeginBlock: func(block *seed.Block) seed.BlockCallbacks {
			blockIdx++
			return seed.BlockCallbacks{
				EndBlock: func() (sealEpoch *pos.Validators) {
					// track blocks
					extended.blocks[blockIdx] = block
					if extended.applyBlock != nil {
						return extended.applyBlock(block)
					}
					return nil
				},
			}
		},
	})
	if err != nil {
		panic(err)
	}

	return extended, store, input
}
