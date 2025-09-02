package keti3

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"keti3/internal/application/helper"

	"keti3/pkg/tal_content_arch_lib/pontos"
)

type TeacherLoginParam struct {
	LoginCode string `json:"login_code" binding:"required"`
}
type TeacherLoginRet struct {
	AuthToken string `json:"authtoken"`
}

func (b business) TeacherLogin(ctx context.Context, param TeacherLoginParam) (ret TeacherLoginRet, err error) {
	loginCodeList := pontos.Config.GetStringSlice("login.teacher_login_code")
	if len(loginCodeList) == 0 {
		err = fmt.Errorf("login.teacher_login_code is empty")
		return
	}
	if !helper.IsInSlice(param.LoginCode, loginCodeList) {
		err = fmt.Errorf("login_code is invalid")
		return
	}
	// 简单生成 token
	hashCode := b.genAuthToken(param.LoginCode)
	// 写入redis
	err = b.saveToken(ctx, hashCode)
	if err != nil {
		return
	}
	ret = TeacherLoginRet{
		AuthToken: hashCode,
	}
	return
}

func (b business) saveToken(ctx context.Context, hashCode string) error {
	hashCode = b.formatRdsKey(hashCode)
	expire := time.Duration(time.Second * 3600)
	// 简单生成 token
	return pontos.Redis.Set(ctx, hashCode, 1, expire).Err()
}

func (b business) formatRdsKey(hashCode string) string {
	preFix := "base-keti:teacher-login:%s"
	return fmt.Sprintf(preFix, hashCode)
}

func (b business) genAuthToken(loginCode string) string {
	hasher := md5.New()
	hasher.Write([]byte(loginCode + tokenSalt))
	hashBytes := hasher.Sum(nil)
	// 将字节数组转换为十六进制字符串表示
	hashString := hex.EncodeToString(hashBytes)
	return hashString
}
