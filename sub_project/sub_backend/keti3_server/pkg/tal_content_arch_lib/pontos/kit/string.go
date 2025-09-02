package kit

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash/crc32"
)

// MD5 ...
func (k *kit) MD5(str string) string {
	c := md5.New()
	c.Write([]byte(str))
	return hex.EncodeToString(c.Sum(nil))
}

// SHA1 ...
func (k *kit) SHA1(str string) string {
	c := sha1.New()
	c.Write([]byte(str))
	return hex.EncodeToString(c.Sum(nil))
}

// SHA256 ...
func (k *kit) SHA256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
}

// CRC32 ...
func (k *kit) CRC32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}
