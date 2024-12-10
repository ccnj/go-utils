package cryptorand

import (
	"crypto/rand"
	"math/big"
)

// 定义可能的字符集
const (
	// 数字和大写字母
	letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// GenCryptoRandStr 生成指定长度的随机字符串
// length: 需要生成的字符串长度
// 返回生成的随机字符串和可能的错误
func GenCryptoRandStr(length int) (string, error) {
	// 创建一个byte切片来存储结果
	result := make([]byte, length)
	// 计算字符集的长度
	lettersLen := big.NewInt(int64(len(letters)))

	for i := 0; i < length; i++ {
		// 为每个位置生成一个随机索引
		num, err := rand.Int(rand.Reader, lettersLen)
		if err != nil {
			return "", err
		}
		// 将对应位置的字符加入结果中
		result[i] = letters[num.Int64()]
	}

	return string(result), nil
}
