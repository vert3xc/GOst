package gost_test

import (
    "encoding/hex"
    "testing"
    "math/big"

    "github.com/vert3xc/GOst/streebog"
    "github.com/vert3xc/GOst/magma"
    "github.com/vert3xc/GOst/kuznyechik"
    "github.com/vert3xc/GOst/gostr34102012"
    "github.com/vert3xc/GOst/gostr1323565.1.006-2017"
)

func TestStreebogHash(t *testing.T) {
    input, _ := hex.DecodeString("323130393837363534333231303938373635343332313039383736353433323130393837363534333231303938373635343332313039383736353433323130")
    expected, _ := hex.DecodeString("486f64c1917879417fef082b3381a4e211c324f074654c38823a7b76f830ad00fa1fbae42b1285c0352f227524bc9ab16254288dd6863dccd5b9f54a1ad0541b")

    h, _ := streebog.New(512)
    h.Write(input)
    sum := h.Sum(nil)

    if !equalBytes(sum, expected) {
        t.Errorf("Streebog hash mismatch\nExpected: %x\nGot:      %x", expected, sum)
    }
}

func TestMagmaEncryptDecrypt(t *testing.T) {
    key, _ := hex.DecodeString("ffeeddccbbaa99887766554433221100f0f1f2f3f4f5f6f7f8f9fafbfcfdfeff")
    plaintext, _ := hex.DecodeString("fedcba9876543210")
    expectedCiphertext, _ := hex.DecodeString("4ee901e5c2d8ca3d")

    c, _ := magma.NewCipher(key)
    ciphertext := make([]byte, 8)
    c.Encrypt(ciphertext, plaintext)

    if !equalBytes(ciphertext, expectedCiphertext) {
        t.Errorf("Magma encrypt mismatch\nExpected: %x\nGot:      %x", expectedCiphertext, ciphertext)
    }

    decrypted := make([]byte, 8)
    c.Decrypt(decrypted, ciphertext)
    if !equalBytes(decrypted, plaintext) {
        t.Errorf("Magma decrypt mismatch\nExpected: %x\nGot:      %x", plaintext, decrypted)
    }
}

func TestKuznyechikEncryptDecrypt(t *testing.T) {
    key, _ := hex.DecodeString("8899aabbccddeeff0011223344556677fedcba98765432100123456789abcdef")
    plaintext, _ := hex.DecodeString("1122334455667700ffeeddccbbaa9988")
    expectedCiphertext, _ := hex.DecodeString("7f679d90bebc24305a468d42b9d4edcd")

    c, _ := kuznyechik.NewCipher(key)
    ciphertext := make([]byte, 16)
    c.Encrypt(ciphertext, plaintext)

    if !equalBytes(ciphertext, expectedCiphertext) {
        t.Errorf("Kuznyechik encrypt mismatch\nExpected: %x\nGot:      %x", expectedCiphertext, ciphertext)
    }

    decrypted := make([]byte, 16)
    c.Decrypt(decrypted, ciphertext)
    if !equalBytes(decrypted, plaintext) {
        t.Errorf("Kuznyechik decrypt mismatch\nExpected: %x\nGot:      %x", plaintext, decrypted)
    }
}

func TestGostSignVerify(t *testing.T) {
    e, _ := new(big.Int).SetString("2DFBC1B372D89A1188C09C52E0EEC61FCE52032AB1022E8E67ECE6672B043EE5", 16)
	k, _ := new(big.Int).SetString("77105C9B20BCD3122823C8CF6FCC7B956DE33814E95B7FE64FED924594DCEAB3", 16)
	privD := new(big.Int)
	privD.SetString("55441196065363246126355624130324183196576709222340016572108097750006097525544", 10)
	curve := gostr34102012.TestParams
	rExpected, _ := new(big.Int).SetString("41AA28D2F1AB148280CD9ED56FEDA41974053554A42767B83AD043FD39DC0493", 16)
	sExpected, _ := new(big.Int).SetString("1456C64BA4642A1653C235A98A60249BCD6D3F746B631DF928014F6C5BF9C40", 16)

    priv := &gostr34102012.GostPrivKey{
        ParentCurve: curve,
        D: privD,
    }
    pub := priv.Public()

    s, r, err := priv.TestSign(e, k)
    if err != nil {
        t.Fatalf("Sign error: %v", err)
    }
	if sExpected.Cmp(s) != 0 || rExpected.Cmp(r) != 0 {
		t.Errorf("Signature computation incorrect\nr: %x\ns: %x\n", r, s)
	}

    ok := gostr34102012.TestVerify(pub, e, r, s)
    if !ok {
        t.Errorf("Signature verification failed")
    }
}

func TestGostRNG(t *testing.T) {
    seed, _ := hex.DecodeString("0101af")
    hsh, _ := streebog.New(512)
    rng := gostr1323565_1_006_2017.New(seed, hsh)

    output := make([]byte, 15)
    _, err := rng.Read(output)
    if err != nil {
        t.Fatalf("RNG read error: %v", err)
    }

    expected, _ := hex.DecodeString("f07d90319a49b637f50d12719be9ec")
    if !equalBytes(output, expected) {
        t.Errorf("RNG output mismatch\nExpected: %x\nGot:      %x", expected, output)
    }
}

func equalBytes(a, b []byte) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}
