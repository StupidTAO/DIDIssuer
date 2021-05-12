package router

import (
	"DIDIssuer/log"
	"DIDIssuer/model"
	"DIDIssuer/utils"
	hub "github.com/StupidTAO/DIDHub/model"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const KEYFILE = "/Users/oker/go/src/github.com/DIDIssuer/private_key"

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "welcome!")
}

func GetClaim(w http.ResponseWriter, r *http.Request) {
	//读取数据
	r.ParseForm()
	rawClaim := r.Form["rawClaim"][0]
	sigClaim := r.Form["sigClaim"][0]
	address := r.Form["addr"][0]

	//签名转换为字节数组
	sigClaimBytes := utils.Base58Decode(sigClaim)

	//验证签名
	b, _ := VerfyClaim(rawClaim, sigClaimBytes, address)
	if b {
		log.Info("signature verify pass")
	}

	//签发声明
	claimId, err := VerifyClaimContent(rawClaim)
	if err != nil {
		log.Info("verify claim content has error: %s", err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}

	//返回声明ID
	fmt.Fprintln(w, claimId)
}

//验证客户端发来的声明数据
func VerfyClaim(rawClaim string, sigClaim []byte, addr string) (bool, error) {
	addrDerive, err := utils.VerifyToAddress(rawClaim, sigClaim)
	if err != nil {
		return false, err
	}
	if addrDerive != addr {
		return false, err
	}

	return true, nil
}

//签发声明并将数据存储到hub
func VerifyClaimContent(rawClaim string) (string, error) {
	result := uint(0)
	csEntry := new(model.CredentialSubject)
	//校验参数
	err := model.CredentialSubjectUnmarshal(rawClaim, csEntry)
	fmt.Println("rawClaim is : ", rawClaim)
	if err != nil {
		return "", err
	}
	if csEntry.ID == "" {
		return "", errors.New("rawClaim parse error")
	}

	if csEntry.TypeCliam == model.IDCardAuthentication {
		//检查身份证号及内容
		_, err := checkIdNumber(csEntry.ShortDescription)
		if err != nil {
			return "", err
		}
		result = 1
	}
	log.Info("claim content is ok")

	bs, err := model.CredentialSubjectMarshal(*csEntry)
	if err != nil {
		return "", err
	}

	pri, err := utils.GetPrivateKeyByFile(KEYFILE)
	if err != nil {
		return "", err
	}
	sigText, err := utils.SignText(string(bs), pri)
	if err != nil {
		return "", err
	}

	proof := new(model.Proof)
	proof.Creator, err = model.GetDIDByPrivateFile(KEYFILE)
	if err != nil {
		return "", err
	}

	proof.ChainAddr = utils.GetAddressByPublicKey(pri.PublicKey)
	proof.EncryptionType = "Keccak256"
	proof.SignatureValue = utils.Base58Encode(sigText)
	log.Info("proof struct is ok")

	//封装proofClaim并存储到didclaim中
	proofClaim := new(model.ProofClaim)
	rand.Seed(time.Now().UnixNano())
	randInt := rand.Intn(10000)
	randStr := strconv.Itoa(randInt)
	id := utils.GetRipemd160HashCode([]byte(randStr))
	proofClaim.Id = utils.Base58Encode(id)
	proofClaim.CredentialSubject = *csEntry
	proofClaim.Proof = *proof
	proofClaim.Result = result
	proofClaim.IssuanceDate = time.Now().Add(8 * time.Hour).String()
	proofClaim.ExpirationDate = time.Now().Add(8 * time.Hour).String()
	proofClaim.Issuer = proof.Creator
	bs, err = model.ProofClaimMarshal(*proofClaim)
	if err != nil {
		return "", err
	}
	log.Info("proof claim struct is ok")

	//封装到数据库数据结构中
	dbClaim := new(hub.DBDIDClaim)
	dbClaim.ClaimId = proofClaim.Id
	dbClaim.DidClaim = string(bs)
	dbClaim.Did = csEntry.ID
	dbClaim.IsAvailable = 1
	dbClaim.CreateTime = time.Now().Add(8 * time.Hour)
	dbClaim.UpdateTime = time.Now().Add(8 * time.Hour)
	//err = model.InsertDBDIDClaim(*dbClaim)
	err = hub.InsertHubDIDClaim(*dbClaim)
	if err != nil {
		return "", err
	}
	log.Info("DBDID claim struct to db is ok")
	return dbClaim.ClaimId, nil
}

//保证年满10周岁
func checkIdNumber(idNumber string) (bool, error) {
	//检查idNumber
	if len(idNumber) != 18 {
		return false, errors.New("id number format error")
	}

	curYear := time.Now().Year()
	year, err := strconv.Atoi(idNumber[6:10])
	if err != nil {
		return false, err
	}
	if curYear - year >= 10 {
		return true, nil
	}
	return false, errors.New("age need bigger than 10 years old")
}
