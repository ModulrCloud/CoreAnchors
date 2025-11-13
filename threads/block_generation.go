package threads

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/modulrcloud/modulr-anchors-core/block_pack"
	"github.com/modulrcloud/modulr-anchors-core/databases"
	"github.com/modulrcloud/modulr-anchors-core/globals"
	"github.com/modulrcloud/modulr-anchors-core/handlers"
	"github.com/modulrcloud/modulr-anchors-core/utils"

	"github.com/syndtr/goleveldb/leveldb"
)

func BlocksGenerationThread() {

	for {

		handlers.APPROVEMENT_THREAD_METADATA.RWMutex.RLock()

		blockTime := handlers.APPROVEMENT_THREAD_METADATA.Handler.NetworkParameters.BlockTime

		generateBlock()

		handlers.APPROVEMENT_THREAD_METADATA.RWMutex.RUnlock()

		time.Sleep(time.Duration(blockTime) * time.Millisecond)

	}

}

func generateBlock() {

	epochHandlerRef := &handlers.APPROVEMENT_THREAD_METADATA.Handler.EpochDataHandler

	if !utils.EpochStillFresh(&handlers.APPROVEMENT_THREAD_METADATA.Handler) {

		return

	}

	epochFullID := epochHandlerRef.Hash + "#" + strconv.Itoa(epochHandlerRef.Id)

	epochIndex := epochHandlerRef.Id

	PROOFS_GRABBER_MUTEX.RLock()

	// Safe "if" branch to prevent unnecessary blocks generation

	shouldGenerateBlocks := handlers.GENERATION_THREAD_METADATA.NextIndex <= PROOFS_GRABBER.AcceptedIndex+1

	shouldRotateEpochOnGenerationThread := handlers.GENERATION_THREAD_METADATA.EpochFullId != epochFullID

	if shouldGenerateBlocks || shouldRotateEpochOnGenerationThread {

		PROOFS_GRABBER_MUTEX.RUnlock()

		// Check if <epochFullID> is the same in APPROVEMENT_THREAD and in GENERATION_THREAD

		if shouldRotateEpochOnGenerationThread {

			// Update the index & hash of epoch (by assigning new epoch full ID)

			handlers.GENERATION_THREAD_METADATA.EpochFullId = epochFullID

			// Nullish the index & hash in generation thread for new epoch

			handlers.GENERATION_THREAD_METADATA.PrevHash = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

			handlers.GENERATION_THREAD_METADATA.NextIndex = 0

		}

		// Safe "if" branch to prevent unnecessary blocks generation
		if !shouldGenerateBlocks {
			return
		}

		extraData := block_pack.ExtraDataToBlock{}

		extraData.Rest = globals.CONFIGURATION.ExtraDataToBlock

		blockDbAtomicBatch := new(leveldb.Batch)

		blockCandidate := block_pack.NewBlock(extraData, epochFullID)

		blockHash := blockCandidate.GetHash()

		blockCandidate.SignBlock()

		// BlockID has the following format => epochID(epochIndex):Ed25519_Pubkey:IndexOfBlockInCurrentEpoch

		blockID := strconv.Itoa(epochIndex) + ":" + globals.CONFIGURATION.PublicKey + ":" + strconv.Itoa(blockCandidate.Index)

		utils.LogWithTime("New block generated "+blockID+" (hash: "+blockHash[:8]+"...)", utils.CYAN_COLOR)

		if blockBytes, serializeErr := json.Marshal(blockCandidate); serializeErr == nil {

			handlers.GENERATION_THREAD_METADATA.PrevHash = blockHash

			handlers.GENERATION_THREAD_METADATA.NextIndex++

			if gtBytes, serializeErr2 := json.Marshal(handlers.GENERATION_THREAD_METADATA); serializeErr2 == nil {

				// Store block locally

				blockDbAtomicBatch.Put([]byte(blockID), blockBytes)

				// Update the GENERATION_THREAD after all

				blockDbAtomicBatch.Put([]byte("GT"), gtBytes)

				if err := databases.BLOCKS.Write(blockDbAtomicBatch, nil); err != nil {

					panic("Can't store GT and block candidate")

				}

			}

		}

	} else {

		PROOFS_GRABBER_MUTEX.RUnlock()

	}

}
