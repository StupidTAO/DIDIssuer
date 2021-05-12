package model

import (
	"DIDIssuer/db"
	"fmt"
	"time"
)

type DBDIDClaim struct {
	Id uint					`db:"id"`
	Did string				`db:"did"`
	ClaimId string			`db:"claimId""`
	DidClaim string			`db:"didClaim"`
	CreateTime time.Time	`db:"createTime"`
	UpdateTime time.Time	`db:"updateTime"`
	IsAvailable	uint		`db:"isAvailable"`
}

func InsertDBDIDClaim(didClaim DBDIDClaim) error {
	sql := "insert into did_claim(did, claimId, didClaim, createTime, updateTime, isAvailable)values (?,?,?,?,?,?)"

	//执行SQL语句
	db.InitDB()
	_, err := db.DB.Exec(sql, didClaim.Did, didClaim.ClaimId, didClaim.DidClaim, didClaim.CreateTime, didClaim.UpdateTime, didClaim.IsAvailable)
	if err != nil {
		fmt.Println("exec failed,", err)
		return err
	}

	return nil
}

func FindDBDIDClaim(did string) ([]DBDIDClaim, error){
	DB := db.InitDB()

	var claims []DBDIDClaim
	sql := "select id, did, claimId, didClaim, createTime, updateTime, isAvailable from did_claim where did=?"
	err := DB.Select(&claims, sql, did)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return nil, err
	}
	return claims, nil
}
