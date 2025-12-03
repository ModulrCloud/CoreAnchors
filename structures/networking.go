package structures

type QuorumMemberData struct {
	PubKey, Url string
}

type AnchorRotationProofRequest struct {
	EpochIndex int        `json:"epochIndex"`
	Creator    string     `json:"creator"`
	Proposal   VotingStat `json:"proposal"`
}

type AnchorRotationProofResponse struct {
	Status     string      `json:"status"`
	Message    string      `json:"message,omitempty"`
	Signature  string      `json:"signature,omitempty"`
	VotingStat *VotingStat `json:"votingStat,omitempty"`
}

type AcceptExtraDataRequest struct {
	RotationProofs []AnchorRotationProof `json:"rotationProofs"`
}

type AcceptExtraDataResponse struct {
	Accepted int `json:"accepted"`
}

type AcceptLeaderFinalizationDataRequest struct {
	LeaderFinalizations []LeaderFinalizationProof `json:"leaderFinalizations"`
}
