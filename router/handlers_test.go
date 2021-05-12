package router

import (
	"DIDIssuer/log"
	"testing"
)

func TestVerifyClaimContent(t *testing.T) {
	//初始化log
	err := log.LogInit()
	if err != nil {
		t.Error("panic: log init error")
		return
	}

	rawClaim := "{\"id\":\"did:welfare:2z7tBiNoYRTCGGNyKcxatEmYxuN1\",\"shortDescription\":\"342225199509082432\",\"longDescription\":\"ID Card\",\"typeClaim\":\"IDCardAuthentication\"} "
	claimId, err := VerifyClaimContent(rawClaim)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(claimId)
}
