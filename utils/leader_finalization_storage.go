package utils

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/modulrcloud/modulr-anchors-core/databases"
	"github.com/modulrcloud/modulr-anchors-core/structures"

	ldbErrors "github.com/syndtr/goleveldb/leveldb/errors"
)

func aggregatedLeaderFinalizationProofKey(epochIndex int, leader string) []byte {

	return []byte("ALFP:" + strconv.Itoa(epochIndex) + ":" + leader)

}

func StoreAggregatedLeaderFinalizationProof(proof structures.AggregatedLeaderFinalizationProof) error {

	payload, err := json.Marshal(proof)

	if err != nil {
		return err
	}

	return databases.FINALIZATION_VOTING_STATS.Put(aggregatedLeaderFinalizationProofKey(proof.EpochIndex, proof.Leader), payload, nil)

}

func LoadAggregatedLeaderFinalizationProof(epochIndex int, leader string) (structures.AggregatedLeaderFinalizationProof, error) {

	var proof structures.AggregatedLeaderFinalizationProof

	raw, err := databases.FINALIZATION_VOTING_STATS.Get(aggregatedLeaderFinalizationProofKey(epochIndex, leader), nil)

	if err != nil {
		if errors.Is(err, ldbErrors.ErrNotFound) {
			return proof, nil
		}
		return proof, err
	}

	if len(raw) == 0 {
		return proof, nil
	}

	if err := json.Unmarshal(raw, &proof); err != nil {
		return proof, err
	}
	return proof, nil

}
