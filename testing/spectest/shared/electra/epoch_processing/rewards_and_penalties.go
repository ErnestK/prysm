package epoch_processing

import (
	"context"
	"path"
	"testing"

	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/electra"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v5/testing/require"
	"github.com/prysmaticlabs/prysm/v5/testing/spectest/utils"
)

// RunRewardsAndPenaltiesTests executes "epoch_processing/rewards_and_penalties" tests.
func RunRewardsAndPenaltiesTests(t *testing.T, config string) {
	require.NoError(t, utils.SetConfig(t, config))

	testPath := "epoch_processing/rewards_and_penalties/pyspec_tests"
	testFolders, testsFolderPath := utils.TestFolders(t, config, "electra", testPath)
	for _, folder := range testFolders {
		helpers.ClearCache()
		t.Run(folder.Name(), func(t *testing.T) {
			folderPath := path.Join(testsFolderPath, folder.Name())
			RunEpochOperationTest(t, folderPath, processRewardsAndPenaltiesPrecomputeWrapper)
		})
	}
}

func processRewardsAndPenaltiesPrecomputeWrapper(t *testing.T, st state.BeaconState) (state.BeaconState, error) {
	ctx := context.Background()
	vp, bp, err := electra.InitializePrecomputeValidators(ctx, st)
	require.NoError(t, err)
	vp, bp, err = electra.ProcessEpochParticipation(ctx, st, bp, vp)
	require.NoError(t, err)

	st, err = electra.ProcessRewardsAndPenaltiesPrecompute(st, bp, vp)
	require.NoError(t, err, "Could not process reward")

	return st, nil
}