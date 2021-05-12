package model

import (
	"DIDIssuer/utils"
	"time"
	"testing"
)

func TestInsertDBDIDClaim(t *testing.T) {
	claim := new(DBDIDClaim)
	claim.Did = "did:welfare:2z7tBiNoYRTCGGNyKcxatEmYxuN1"
	claim.DidClaim = "claim-test"
	bs := utils.GetRipemd160HashCode([]byte(claim.DidClaim))
	claim.ClaimId = utils.Base58Encode(bs)
	claim.IsAvailable = 1
	claim.CreateTime = time.Now().Add(8 * time.Hour)
	claim.UpdateTime = time.Now().Add(8 * time.Hour)
	err := InsertDBDIDClaim(*claim)
	if err != nil {
		t.Error("error is ", err)
	}
}

func TestFindDBDIDClaim(t *testing.T) {
	did := "did:welfare:2z7tBiNoYRTCGGNyKcxatEmYxuN1"
	claims, err := FindDBDIDClaim(did)
	if err != nil {
		t.Error("find db did claims error ", err)
		return
	}
	t.Log("claims counts : ", len(claims))
}
