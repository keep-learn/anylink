package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"os"

	"github.com/bjdgyc/anylink/admin"
	"github.com/bjdgyc/anylink/base"
	"github.com/bjdgyc/anylink/dbdata"
	"github.com/bjdgyc/anylink/sessdata"
)

func Start() {
	dbdata.Start()
	sessdata.Start()

	switch base.Cfg.LinkMode {
	case base.LinkModeTUN:
		checkTun()
	case base.LinkModeTAP:
		checkTap()
	case base.LinkModeMacvtap:
		checkMacvtap()
	default:
		base.Fatal("LinkMode is err")
	}

	// 计算profile.xml的hash
	b, err := os.ReadFile(base.Cfg.Profile)
	if err != nil {
		panic(err)
	}
	ha := sha1.Sum(b)
	profileHash = hex.EncodeToString(ha[:])

	// 主要是后端管理系统的一些路由
	go admin.StartAdmin()

	// 这个是tcp版本的加密
	go startTls()

	// 这个是udp版本的加密
	go startDtls()
}

func Stop() {
	_ = dbdata.Stop()
	destroyVtap()
}
