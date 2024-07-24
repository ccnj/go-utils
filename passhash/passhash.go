package passhash

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

// GenerateSalt generates a random salt of specified length
func GenerateSalt(length int) ([]byte, error) {

	// alpha字母 num数字
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		// 错误与底层操作系统或硬件提供的随机数生成器有关，概率很小
		return nil, err
	}
	for key, val := range salt {
		// 与长度取余数，保证每个字符都是alphanum中的字符。
		// 但并不是完全随机落在每个字符上，因为255（val为0-255）不一定是len(alphanum)的整数倍。
		// 不过也基本上是均匀分布的。
		// 不直接用val是因为，最终这个盐要以string形式存在数据库里，val为0-255（ASCII码），转string有的是乱码
		salt[key] = alphanum[val%byte(len(alphanum))]
	}
	return salt, nil
}

// HashPassword hashes a password using PBKDF2 with a given salt and iteration count.
// keyLen is the length of the resulting hash in bytes，可任意改 33 34 64 啥都行
// iterations is the number of iterations to use in PBKDF2 , 既重复进行多少次哈希计算(对哈希后的结果再哈希，就更难破译)
func HashPassword(rawPwd string, salt []byte) string {
	// 写死了，迭代100次，hash长度32字节，sha512加密算法
	hashBytes := pbkdf2.Key([]byte(rawPwd), salt, 100, 32, sha512.New)

	// hashBytes是[]byte类型，不易存储于数据库，所以要编码
	// 1. 转成hex字符串，方便存储，但长度的len(hashBytes)的两倍，因为两个hex（0-15）才能存一个byte（0-255）
	// return hex.EncodeToString(hashBytes)
	// 2. 转为base64，base64就是为了把二进制数据（一个byte就是4位二进制数据）编码成可打印字符
	// 每三个字节转换为四个 Base64 字符，大约是len(hashBytes)的1.33倍
	return base64.StdEncoding.EncodeToString(hashBytes)
}

// 传入原始密码和盐长度，返回盐和哈希后的密码
// err只有在生成盐时才可能出错，与底层操作系统或硬件提供的随机数生成器有关，概率很小
func EasyHash(rawPwd string, saltLen int) (salt string, pwdHash string, err error) {
	saltBytes, err := GenerateSalt(saltLen)
	if err != nil {
		return "", "", err
	}
	pwdHash = HashPassword(rawPwd, saltBytes)
	return string(saltBytes), pwdHash, nil
}

// 验证密码是否正确，传入原始密码、哈希后的密码和盐
func Verify(rawPwd, pwdHash string, salt string) bool {
	hash := pbkdf2.Key([]byte(rawPwd), []byte(salt), 100, 32, sha512.New)
	return base64.StdEncoding.EncodeToString(hash) == pwdHash
}
