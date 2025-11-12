package structures

type ApprovementThreadMetadataHandler struct {
	NetworkParameters NetworkParameters `json:"networkParameters"`
	EpochDataHandler  EpochDataHandler  `json:"epoch"`
}

func (handler *ApprovementThreadMetadataHandler) GetNetworkParams() NetworkParameters {
	return handler.NetworkParameters
}

func (handler *ApprovementThreadMetadataHandler) GetEpochHandler() EpochDataHandler {
	return handler.EpochDataHandler
}

type GenerationThreadMetadataHandler struct {
	EpochFullId string `json:"epochFullId"`
	PrevHash    string `json:"prevHash"`
	NextIndex   int    `json:"nextIndex"`
}
