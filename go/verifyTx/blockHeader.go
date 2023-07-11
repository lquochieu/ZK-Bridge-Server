package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	//"io/ioutil"
	//"os"
	//"strconv"
	"time"

	//"github.com/ethereum/go-ethereum/common/hexutil"
	//"github.com/tendermint/tendermint/crypto/ed25519"
	//"github.com/tendermint/tendermint/light"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
)

type SignedHeader struct {
	Header Header `json:"header"`
	Commit Commit `json:"commit"`
}

type PartSetHeader struct {
	Total uint32 `json:"total"`
	Hash string `json:"hash"`
}

type Version struct {
	Block string `json:"block"`
	App PartSetHeader `json:"app"`
}

type Block_ID struct {
	Hash string `json:"hash"`
	Parts PartSetHeader `json:"parts"`

}
type Header struct {
	Version Version `json:"version"`
	ChainID string `json:"chain_id"`
	Height string `json:"height"`
	Time string `json:"time"`
	LastBlockID Block_ID `json:"last_block_id"`
	LastCommitHash string `json:"last_commit_hash"`
	DataHash string `json:"data_hash"`
	ValidatorsHash string `json:"validators_hash"`
	NextValidatorsHash string `json:"next_validators_hash"`
	ConsensusHash string `json:"consensus_hash"`
	AppHash string `json:"app_hash"`
	LastResultsHash string `json:"last_results_hash"`
	EvidenceHash string `json:"evidence_hash"`
	ProposerAddress string `json:"proposer_address"`
}

func GetBlockHeader(signed_header_path string) types.SignedHeader {
	var signed_header SignedHeader
	json.Unmarshal(GetBytesFileJson(signed_header_path), &signed_header)
	var commit_sigs []types.CommitSig
	
	// commit_sigs_pointer := make([]*types.CommitSig, len(signed_header.Commit.Signatures))
    for i := 0; i < len(signed_header.Commit.Signatures); i++ {
		commit_sigs = append(commit_sigs, types.CommitSig {
			BlockIDFlag: types.BlockIDFlag(signed_header.Commit.Signatures[i].BlockIDFlag),
			ValidatorAddress: First(hex.DecodeString(signed_header.Commit.Signatures[i].ValidatorAddress)),
			Timestamp:        First(time.Parse(time.RFC3339, signed_header.Commit.Signatures[i].Timestamp)),
			Signature:        First(base64.StdEncoding.DecodeString(signed_header.Commit.Signatures[i].Signature)),
		})
		// fmt.Println("address:", hexutil.Encode(First(hex.DecodeString(signed_header.Commit.Signatures[i].ValidatorAddress))))
		// fmt.Println("signature",First(base64.StdEncoding.DecodeString(signed_header.Commit.Signatures[i].Signature)))
		// fmt.Println(i, First(time.Parse(time.RFC3339, signed_header.Commit.Signatures[i].Timestamp)));
		// commit_sigs_pointer[i] =  &commit_sigs[i]
	}
	fmt.Println("")
	var header = types.SignedHeader{
		Header: &types.Header{
			Version: version.Consensus{
				Block: uint64(StringToInt64(signed_header.Header.Version.Block)),
			},
			ChainID: signed_header.Header.ChainID,
			Height:  StringToInt64(signed_header.Header.Height),
			Time:    First(time.Parse(time.RFC3339, signed_header.Header.Time)),
			LastBlockID: types.BlockID{
				Hash: First(hex.DecodeString(signed_header.Header.LastBlockID.Hash)),
				PartSetHeader: types.PartSetHeader{
					Total: signed_header.Header.LastBlockID.Parts.Total,
					Hash:  First(hex.DecodeString(signed_header.Header.LastBlockID.Parts.Hash)),
				},
			},
			LastCommitHash:     First(hex.DecodeString(signed_header.Header.LastCommitHash)),
			DataHash:           First(hex.DecodeString(signed_header.Header.DataHash)),
			ValidatorsHash:     First(hex.DecodeString(signed_header.Header.ValidatorsHash)),
			NextValidatorsHash: First(hex.DecodeString(signed_header.Header.NextValidatorsHash)),
			ConsensusHash:      First(hex.DecodeString(signed_header.Header.ConsensusHash)),
			AppHash:            First(hex.DecodeString(signed_header.Header.AppHash)),
			LastResultsHash:    First(hex.DecodeString(signed_header.Header.LastResultsHash)),
			EvidenceHash:       First(hex.DecodeString(signed_header.Header.EvidenceHash)),
			ProposerAddress:    First(hex.DecodeString(signed_header.Header.ProposerAddress)),
		},
		Commit: &types.Commit{
			Height: StringToInt64(signed_header.Commit.Height),
			Round:  signed_header.Commit.Round,
			BlockID: types.BlockID{
				Hash: First(hex.DecodeString(signed_header.Commit.BlockID.Hash)),
				PartSetHeader: types.PartSetHeader{
					Total: uint32(signed_header.Commit.BlockID.Parts.Total),
					Hash:  First(hex.DecodeString(signed_header.Commit.BlockID.Parts.Hash)),
				},
			},
			Signatures: commit_sigs,
		},
	}
	// fmt.Println(First(time.Parse(time.RFC3339, signed_header.Header.Time)));
	return header
}