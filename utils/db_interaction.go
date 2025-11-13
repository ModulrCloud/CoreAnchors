package utils

import (
	"encoding/json"

	"github.com/modulrcloud/modulr-anchors-core/databases"
	"github.com/modulrcloud/modulr-anchors-core/globals"
	"github.com/modulrcloud/modulr-anchors-core/handlers"
	"github.com/modulrcloud/modulr-anchors-core/structures"

	"github.com/syndtr/goleveldb/leveldb"
)

func OpenDb(dbName string) *leveldb.DB {

	db, err := leveldb.OpenFile(globals.CHAINDATA_PATH+"/DATABASES/"+dbName, nil)
	if err != nil {
		panic("Impossible to open db : " + dbName + " =>" + err.Error())
	}
	return db
}

func GetAnchorFromApprovementThreadState(validatorPubkey string) *structures.AnchorsStorage {

	validatorStorageKey := validatorPubkey + "_ANCHOR_STORAGE"

	if val, ok := handlers.APPROVEMENT_THREAD_METADATA.Handler.ValidatorsStoragesCache[validatorStorageKey]; ok {
		return val
	}

	data, err := databases.APPROVEMENT_THREAD_METADATA.Get([]byte(validatorStorageKey), nil)

	if err != nil {
		return nil
	}

	var validatorStorage structures.AnchorsStorage

	err = json.Unmarshal(data, &validatorStorage)

	if err != nil {
		return nil
	}

	handlers.APPROVEMENT_THREAD_METADATA.Handler.ValidatorsStoragesCache[validatorStorageKey] = &validatorStorage

	return &validatorStorage

}
