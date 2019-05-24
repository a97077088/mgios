package main

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
)

//PKCS7去除填充
func PKCS7UnPadding_with_in(_in []byte) []byte {
	length := len(_in)
	unpadding := int(_in[length-1])
	return _in[:length-unpadding]
}

//使用PKCS7进行填充
func PKCS7Padding_with_in_blocksize(_in []byte, _blockSize int) []byte {
	padding := _blockSize - len(_in)%_blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(_in, padtext...)
}

//使用PKCS5进行填充
func PKCS5Padding_with_in_blocksize(_in []byte, _blockSize int) []byte {
	padding := _blockSize - len(_in)%_blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(_in, padtext...)
}

func Hmac_sha256_with_in_key(_in []byte, _key []byte) []byte {
	h := hmac.New(sha256.New, _key)
	h.Write(_in)
	return h.Sum(nil)
}
func Sha256_with_in(_in []byte) []byte {
	h := sha256.New()
	h.Write(_in)
	return h.Sum(nil)
}
func Rsa_with_in_privatekey(_in []byte, _privatekey []byte) ([]byte, error) {
	block, _ := pem.Decode(_privatekey)
	privkey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	var (
		ErrInputSize  = errors.New("input size too large")
		ErrEncryption = errors.New("encryption error")
	)
	k := (privkey.N.BitLen() + 7) / 8
	tLen := len(_in)
	// rfc2313, section 8:
	// The length of the data D shall not be more than k-11 octets
	if tLen > k-11 {
		return nil, ErrInputSize
	}
	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k-tLen-1; i++ {
		em[i] = 0xff
	}
	copy(em[k-tLen:k], _in)
	c := new(big.Int).SetBytes(em)
	if c.Cmp(privkey.N) > 0 {
		return nil, ErrEncryption
	}
	var m *big.Int
	var ir *big.Int
	if privkey.Precomputed.Dp == nil {
		m = new(big.Int).Exp(c, privkey.D, privkey.N)
	} else {
		// We have the precalculated values needed for the CRT.
		m = new(big.Int).Exp(c, privkey.Precomputed.Dp, privkey.Primes[0])
		m2 := new(big.Int).Exp(c, privkey.Precomputed.Dq, privkey.Primes[1])
		m.Sub(m, m2)
		if m.Sign() < 0 {
			m.Add(m, privkey.Primes[0])
		}
		m.Mul(m, privkey.Precomputed.Qinv)
		m.Mod(m, privkey.Primes[0])
		m.Mul(m, privkey.Primes[1])
		m.Add(m, m2)

		for i, values := range privkey.Precomputed.CRTValues {
			prime := privkey.Primes[2+i]
			m2.Exp(c, values.Exp, prime)
			m2.Sub(m2, m)
			m2.Mul(m2, values.Coeff)
			m2.Mod(m2, prime)
			if m2.Sign() < 0 {
				m2.Add(m2, prime)
			}
			m2.Mul(m2, values.R)
			m.Add(m, m2)
		}
	}

	if ir != nil {
		// Unblind.
		m.Mul(m, ir)
		m.Mod(m, privkey.N)
	}
	enc := m.Bytes()
	return enc, nil
}
func Must_Rsa_with_in_privatekey(_in []byte, _privatekey []byte) []byte {
	r, err := Rsa_with_in_privatekey(_in, _privatekey)
	if err != nil {
		panic(err)
	}
	return r
}
func Rsa_with_in_publickey(_in []byte, _publickey []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode([]byte(_publickey))
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, _in)
}
func Rsa_decode_with_in_privatekey(_in []byte, _privatekey []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(_privatekey)
	if block == nil {
		return nil, errors.New("解析pem失败")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	outbuf, err := rsa.DecryptPKCS1v15(rand.Reader, priv, _in)
	if err != nil {
		return nil, err
	}
	return outbuf, nil
}
func Rsa_sign_md5_with_in_privatekey(_in []byte, _privatekey []byte) ([]byte, error) {
	block, _ := pem.Decode(_privatekey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	h := crypto.MD5.New()
	h.Write(_in)
	return rsa.SignPKCS1v15(rand.Reader, priv.(*rsa.PrivateKey), crypto.MD5, h.Sum(nil))
}
func Must_Rsa_sign_md5_with_in_privatekey(_in []byte, _privatekey []byte) []byte {
	r, err := Rsa_sign_md5_with_in_privatekey(_in, _privatekey)
	if err != nil {
		panic(err)
	}
	return r
}
func Rsa_sign_sha256_with_in_privatekey(_in []byte, _privatekey []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode([]byte(_privatekey))
	if block == nil {
		return nil, errors.New("private key error")
	}
	// 解析公钥
	privatekey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rng := rand.Reader
	hashed := sha256.Sum256(_in)
	return rsa.SignPKCS1v15(rng, privatekey, crypto.SHA256, hashed[:])
}
func Must_Rsa_sign_sha256_with_in_privatekey(_in []byte, _privatekey []byte) []byte {
	r, err := Rsa_sign_sha256_with_in_privatekey(_in, _privatekey)
	if err != nil {
		panic(err)
	}
	return r
}
func Must_Md5_with_in(_in []byte) []byte {
	h := crypto.MD5.New()
	h.Write(_in)
	return h.Sum(nil)
}
func Aes128_with_in_key(_in []byte, _key []byte) ([]byte, error) {
	block, err := aes.NewCipher(_key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData := PKCS7Padding_with_in_blocksize(_in, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}
func Must_Aes128_with_in_key(_in []byte, _key []byte) []byte {
	r, err := Aes128_with_in_key(_in, _key)
	if err != nil {
		panic(err)
	}
	return r
}
func Aes128_decode_with_in_key(_in []byte, _key []byte) ([]byte, error) {
	block, _ := aes.NewCipher(_key)
	blockMode := cipher.NewCBCDecrypter(block, []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
	origData := make([]byte, len(_in))
	blockMode.CryptBlocks(origData, _in)
	origData = PKCS7UnPadding_with_in(origData)
	return origData, nil
}

var (
	MGStatic_rsaString_privatekey = `-----BEGIN RSA PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBALrItgh36kGPNVLf
S19vLv/fCS1aSuezWm9lmjvL/fTbAn0tCkQiFAztQ794z4GUBc1HpYy0+RYUzpYk
XwLKLBccZtFTYSJC/cDSaNKc9DWEwjhi1UZCkjZIqQT6pHG+bhQECPExQJZNyfWH
++kijbiVXDsCm7+DmKkowfkCQJJdAgMBAAECgYABdNUuhbjW4AX7anBNAzHrntKU
fvK8kcJvwewM8U8+nIKNdU9yLo3Xd/YiRQ/hI8VTz240Bak3idv+gbQX8nmbhRh2
DmLFEssQK2JAocxviYb4vSnR6kl8aQ4JkBTFK5sD1DggEf08PeDQRG2VVcJHba57
w84pSAJJycc98eNF4QJBAPKr4gVwuiw0xOSbzE6J23YTafP58egfK6ntclxvXCBb
wO24/dkHvZSKo65KXmedyn9J+eDF2eh60XzcbT4iLtUCQQDFCwPr1Z/qoWztFqkZ
Rx11IOVpEOUVHQZXXpUQQ3yhfFB3Bilkl+qztXpO/0/xvLyCU6iwpsL0AWNpQFyb
7mlpAkAE/NlpaTa6X99IH5ekaZr0I1QDKrj5H2LjXDkEg1luZNRjTnERO+J+ohbZ
Jzgngwu9cATNxONErvmNSm6IX2V9AkBQmcjo12Gq4I1mDjbtmCF8oBA3gimzEURe
yNnA5Vo91FQ129dixqHDESEO33EcOVADH6WJXS2yEj9P9nvrbQLBAkBvZTOsP6WG
T9UMv4bSrL0n/mvC2S/xfruTIY3RJ4pKWhnKYH4SiQErfv9fE9WA13r1idLhk90a
5JIcxTJO8TQo
-----END RSA PRIVATE KEY-----`
	//初始化密钥,定位特征
	//ApiHandler requestPlayURLWithParams:success:failure:
	//MGTEncryptWrapper getEncryptWithTimestamp:
	//进两层off_指向
	GetMiguSign_arrindex = []string{
		"57e7f108bfe849d1b7cdfe8ee09d253c",
		"e20c2f9c8007461b823b568f7e2f1df8",
		"de6d4d4e764f4bd9b9d3c61482dc0e19",
		"aeba4e6604884a3aad4c9c1f8e19927d",
		"6dc9ecf9c4414f1a8ba43f340105bc5b",
		"1d88a1f67977431f8a0bd3e37ed458a0",
		"0b5b7a1418d040faad7cbfa029c91ce5",
		"42984878daf74ca7b827ededee590953",
		"117845d9550446bfa4f2c158cdd1fd73",
		"e3eb9ae03571433a99ea5099bce0fc80",
		"7fdb3be24a8a472982c871297f06cd84",
		"0dc13af15fcb4f59bf6fa73979c3922b",
		"11e72b06381d4e11b2d09138c5f438a7",
		"d63cbbafdadc468eaf08ed41eeb9a005",
		"882578c312224526b92ccf4ce8adbbeb",
		"1d0f73d53f4444488884bffb33632076",
		"17261d6d179940b48a0bcaed5d866dfe",
		"ccfb0cac55ef4a66bf1570bdc993a112",
		"49ca01d63676425a8a4da851237c3746",
		"e925ec0094534595995c7cc0ee5d497e",
		"8e1aa9aade3f44df985ed674df81636f",
		"7f163c264789422fbb4e94601e57b5d2",
		"7a0708a4a1514858a70d2925f9e652f6",
		"fd570de139d94a82af8534fd718e77f5",
		"f18a099403394af5b6704edebdefc808",
		"531010e9b1f84023890a5b117abdf5cd",
		"b85de9ea18a94270ab000a6bf470f3d2",
		"c482d9f60b8240909caa1ff93eaa614a",
		"b96a5c3c506c41079fb44974dab80692",
		"af1c58d2c4ef49d19c2dbc8999bc689a",
		"da22fc905b9c4253abfb896725b45563",
		"1c88a1a93b1d42ed87a67f6aa515d4c9",
		"7c3fa7851f21495d844d4d3ed10b7b3e",
		"cfc38b72bd0244839ccdc144f7884251",
		"8a5431431581404f906c8fbb61bb509f",
		"0114c865b5b143c69e665297d09cde21",
		"da0ff36022c94d82ae32b861c6203d0a",
		"deb557c0976c48c59fd2a8d597d92e78",
		"d5965cd090184a78a60ef4805055fbd9",
		"27a1a6c25862493ebec83e1fcb8783bd",
		"5e73077f1ce14d7cac6099d97f360c1d",
		"d52695923cc4467abc52e354560a7cb0",
		"a71e25f15165489ba16c34c2e45c435e",
		"41b228fc1acd4aa68f7669bfefaeceb8",
		"c29c6e07bc324670a6404362c399f7d4",
		"8da92c6775564d7f87ee1df9e367656b",
		"1c35312f70ed4291a686e1ffb11b0426",
		"07a423e0b15d43bdaadb0d538bdd0bda",
		"403c8c4730d042ca965bebd81d079153",
		"1d068f637c1a479886922ce6763a0c4f",
		"ba569bf27cd24a6faa3fc3478c8e87e4",
		"d02e5299c8d243b3a261abe72b979ab3",
		"4103d0198b5d4fef9985a5bd34d8122a",
		"8705d4615a824c9790d4d25792bfdc82",
		"73f42fe46c204b05a0d5222bc3e433fb",
		"ba78f5e849c240e48b60f9e49f8dd1be",
		"a6acc8e4942248c89cdd8ff4cd1b8aa0",
		"af69737f912a43c984b7c4f5ddbe7c99",
		"92eaec15b90e4c05a76bdce7f3bd2637",
		"06fb6827b8244cf495620699e1aa6eb1",
		"05c4e01d1e904cc18609c1ee3fe43a53",
		"9070d9eb650841c6b14f9868593ddd16",
		"4c8065d58f424ff6ab279c53557f1030",
		"1b3b1143d5ef472596d5988523229ee3",
		"d97918056fee43d9a1eaa02fd0d3b1cf",
		"82eafa984900413d91bd45472853c972",
		"7f01eab066754c689fdd05e46e6d825c",
		"34790025cb6f46bcae3fe7f9bd5a8489",
		"2f5011fb2a2648e5883774920c373a34",
		"72312fbeb2e444f8941fd650678ab4b6",
		"55f7860c0af04f29bc6f8a27762dd201",
		"4c963cb27bec45d3a81a47ae0bbe1e8e",
		"6d83f57cee6c408f819b8a7c89985b26",
		"56457f86ec2f43b0a6ff8de6e3d1bb05",
		"ffd8bd316a6948eca113bd232961c1c3",
		"315a49e821ee4128b992716cea3f65c3",
		"8a64cac647c5443ebc3c081e64695587",
		"841c9d63a887435fb6264a388fc2ca28",
		"0cd69540164c4c689299166b7774952b",
		"7ef85021241745a288641914587befa5",
		"e494dd0e4f4f42b49f027d3cc87b9d8f",
		"b7a2484c210540fa951c0e74e8f4abda",
		"eeaf6717c7474349afb76f8051ac7466",
		"342634215c904533827ca9558776622c",
		"0c3680ca49ad4d91a837f413a9ae6147",
		"8524acd7e2b64137a9d218178b63528b",
		"25be331bcc034a53bb698061b293844a",
		"b41d563d48ce483aae4f84491f30f910",
		"d5fe1432988e4dd5a0a0aad27f47e786",
		"2b5edde5b2cc4ea88be381647f5d38a4",
		"dee6a59c81ec4130af1652db14efdb77",
		"efea4dc7a8e4460a995fa55f19f73227",
		"b24e6adb29b346d8a3022d596ec9392f",
		"6a91a2649bd142cb89fd74bc45c04684",
		"db7014ddd71f4de58f1c308e6e79a312",
		"d120e89a5fc74228b70bb10e192ae325",
		"712ede6028b640d6b92f33d5694e089a",
		"77177921d86e4d4394bac284e7baef3b",
		"0001bce083b348d6a50aa1a1149d908f",
		"0e41c2b6b01a46b6952affb9f78ac2b8"}
	AuthRequestpubkey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDMQK5vjD6+FRf8/RQDgeyg8T/G
vLo3CD7ngcZ0l2fXuHh83gxgM1C7BNN1RgQFRXWFzMDLxbynm/ZA3GsntY+OiN8S
Bza2SEnir+p0jrsZ+WNhaRWch4FpAxnfcnfxrJL26ZkfYaxS85pKby/IQuSwIswV
om0F06RIZreK3rQDrwIDAQAB
-----END PUBLIC KEY-----`
	Signature_privatekey = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQDN2CQ+ZcsbHVrLSJ+/v0BA6kgoa+riowfZErlGdlMjUhf4UTos
1ZzKfavQjctIQjlEStxBWAZ+gtmtYakUmjWFBupks6RWcR/ho+VEePpzxhSis8M6
zHDhf9Dl8u6uvSAdk+CM3NgGle0bC8WnMxy2k9k3MSrX93AKZV36DCMtUQIDAQAB
AoGACBp+MsanHEYnkOEnCNFqoiOW+6Bj+tAYOv91s8RsuXM95lSsSZ+PMJmJ7gfm
/M0+m+Wmjhv9BXX5Q84Ybes0OBpS9qk2Rq6mQooXqo+6BaPlJb/UD160ULiQJIA7
P2x2XX/Z8xW7goq6r13i4VOZj4GHxRnlfvtCyKKso2U4qMECQQD2bBrY+1O1Du7D
B6jjJC0rnZMQllIlex9zftvIgr69qvQcWoCKukBUjmx8gw7Pall33sRgJbhMj0yf
qLBn6GplAkEA1dhJcJhZ9euFsjWnjXG5g5CyRO5HPJ9Kfoqw04aQX7oHHRXUyVNU
H4L9FsbXQBYS+rXrtJUcas1Ns7inhamyfQJAXwvndxXJfaaa1ULZE3NasN4AYX95
g9dvlB60Kyyy4XlU2rLVrayVL4gXtBbg2YPNqnyUBfnGklEbXuGz2QA+OQJAUOiQ
gMVj5CPEZfTe/Ck3I4wvptzwnwM10ELxPFcBcPaVkm+cHsAkZ/fLgj6hWmH/tFP4
Zk60fcRHzePjnjLikQJAcaztmnB32RZNP3T0EWDT/ByFKjUl53epXX8YIfxipa1M
sWH66CeJ3S1o8G3CORXOPSJbhDCRBQnWLH+Dp3LovQ==
-----END RSA PRIVATE KEY-----`
	Login_migutokenforall_rsasign_privatekey = `-----BEGIN PRIVATE KEY-----
MIICXAIBAAKBgQDBFOnukK/bomnQB6rSupd7Cgh71Pwbylgbm2kuWRmH25BadjqggMsoU/i4uycsjj
oG6CwqXugeq1JnUV7hF6jAiWz7FsYrtESPAHnOq40C7vpSSrPt1q7oHMnpLReYVHrbcDc+uEhE5Bbrad
C38c32S5jxjpodn87v/5zLeCFqqwIDAQABAoGAHiRcrLCLs5b0O9Sml0Un1r5nOqWyQchh3tVxIxonwS
zGqUihuGLC1GXfgz88S1lct61RD8BHLlqCf7yVOkCOSFMBSKSg9VBFrZb+7ogWxmGPpjDL4yRehwuLRL
U/kdLOjLN0tub/bnNtZ/z9zP/SihVQmLL/dPbPCEMJmFmZOBECQQDfLWkoDYylpRbR7RG8tKysDbNyLW
oQI+2zSGSY40+jv/tAl4dk2JUA3yjx0ynPF83gX7Ogabz/g5ZyN0lv6CQDAkEA3Xpn0ps/N5WKB9pLAe
EO0LK5zvCo2gsUlVYsadwuAerBTwqGqlRNQSEBFrw5GGOWcbZZ/X+ff/Dkv3gW7uMiOQJBAJzHkuK/Of
m3kuNgLiCpr8+iRFhGTQcPplFW3syEixLWsBN9H3EZsPsOyf2vwOQprcgcktxyy4GYGB0ed6l0o68CP1
TIhvYYeBrPJAfdiNkvTnyV7otJlVni4/5G/rmHkBUryNR5MxQBMZG9EK2jYDf156GgLNLeCDDyBp9FWY
IZcQJBANRxXgq6IUXKpgvTFCHeTCHoNztOvN4XS8MYWkZK34a/1/mGLs9U+K1lnBnPho+ZtEpYyW5VZk
awSHjDplzlM8U=
-----END PRIVATE KEY-----`
	Login_migutokenforall_uadecode_privatekey = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQDN2CQ+ZcsbHVrLSJ+/v0BA6kgoa+riowfZErlGdlMjUhf4UTos
1ZzKfavQjctIQjlEStxBWAZ+gtmtYakUmjWFBupks6RWcR/ho+VEePpzxhSis8M6
zHDhf9Dl8u6uvSAdk+CM3NgGle0bC8WnMxy2k9k3MSrX93AKZV36DCMtUQIDAQAB
AoGACBp+MsanHEYnkOEnCNFqoiOW+6Bj+tAYOv91s8RsuXM95lSsSZ+PMJmJ7gfm
/M0+m+Wmjhv9BXX5Q84Ybes0OBpS9qk2Rq6mQooXqo+6BaPlJb/UD160ULiQJIA7
P2x2XX/Z8xW7goq6r13i4VOZj4GHxRnlfvtCyKKso2U4qMECQQD2bBrY+1O1Du7D
B6jjJC0rnZMQllIlex9zftvIgr69qvQcWoCKukBUjmx8gw7Pall33sRgJbhMj0yf
qLBn6GplAkEA1dhJcJhZ9euFsjWnjXG5g5CyRO5HPJ9Kfoqw04aQX7oHHRXUyVNU
H4L9FsbXQBYS+rXrtJUcas1Ns7inhamyfQJAXwvndxXJfaaa1ULZE3NasN4AYX95
g9dvlB60Kyyy4XlU2rLVrayVL4gXtBbg2YPNqnyUBfnGklEbXuGz2QA+OQJAUOiQ
gMVj5CPEZfTe/Ck3I4wvptzwnwM10ELxPFcBcPaVkm+cHsAkZ/fLgj6hWmH/tFP4
Zk60fcRHzePjnjLikQJAcaztmnB32RZNP3T0EWDT/ByFKjUl53epXX8YIfxipa1M
sWH66CeJ3S1o8G3CORXOPSJbhDCRBQnWLH+Dp3LovQ==
-----END RSA PRIVATE KEY-----`
)
