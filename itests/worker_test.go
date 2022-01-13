package itests

import (
	"context"
	"github.com/filecoin-project/lotus/build"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/lotus/extern/sector-storage/sealtasks"
	"github.com/filecoin-project/lotus/itests/kit"
)

func TestWorkerPledge(t *testing.T) {
	ctx := context.Background()
	_, miner, worker, ens := kit.EnsembleWorker(t, kit.WithAllSubsystems(), kit.ThroughRPC(), kit.WithNoLocalSealing(true),
		kit.WithTaskTypes([]sealtasks.TaskType{sealtasks.TTFetch, sealtasks.TTCommit1, sealtasks.TTFinalize, sealtasks.TTAddPiece, sealtasks.TTPreCommit1, sealtasks.TTPreCommit2, sealtasks.TTCommit2, sealtasks.TTUnseal})) // no mock proofs

	ens.InterconnectAll().BeginMining(50 * time.Millisecond)

	e, err := worker.Enabled(ctx)
	require.NoError(t, err)
	require.True(t, e)

	miner.PledgeSectors(ctx, 1, 0, nil)
}

func TestWinningPostWorker(t *testing.T) {
	prevIns := build.InsecurePoStValidation
	build.InsecurePoStValidation = false
	defer func() {
		build.InsecurePoStValidation = prevIns
	}()

	ctx := context.Background()
	client, _, worker, ens := kit.EnsembleWorker(t, kit.WithAllSubsystems(), kit.ThroughRPC(),
		kit.WithTaskTypes([]sealtasks.TaskType{sealtasks.TTGenerateWinningPoSt})) // no mock proofs

	ens.InterconnectAll().BeginMining(50 * time.Millisecond)

	e, err := worker.Enabled(ctx)
	require.NoError(t, err)
	require.True(t, e)

	client.WaitTillChain(ctx, kit.HeightAtLeast(6))
}
