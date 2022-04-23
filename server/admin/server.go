// admin:后台管理接口
package admin

import (
	"crypto/tls"
	"embed"
	"net/http"
	"net/http/pprof"

	"github.com/bjdgyc/anylink/base"
	"github.com/gorilla/mux"
)

var UiData embed.FS

// StartAdmin 开启服务 主要是管理系统的路由信息
func StartAdmin() {

	// 使用 mux 框架：这个框架和gin类似
	r := mux.NewRouter()
	// mux 中也有类似gin的middleware ，长见识了
	r.Use(authMiddleware)

	// 监控检测
	r.HandleFunc("/status.html", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Name("index")

	// 前后端在一个项目中，开发方便；后期做大了，以前后端分离
	r.Handle("/", http.RedirectHandler("/ui/", http.StatusFound)).Name("index")
	r.PathPrefix("/ui/").Handler(
		// http.StripPrefix("/ui/", http.FileServer(http.Dir(base.Cfg.UiPath))),
		http.FileServer(http.FS(UiData)),
	).Name("static")
	r.HandleFunc("/base/login", Login).Name("login")

	r.HandleFunc("/set/home", SetHome)
	r.HandleFunc("/set/system", SetSystem)
	r.HandleFunc("/set/soft", SetSoft)
	r.HandleFunc("/set/other", SetOther)
	r.HandleFunc("/set/other/edit", SetOtherEdit)
	r.HandleFunc("/set/other/smtp", SetOtherSmtp)
	r.HandleFunc("/set/other/smtp/edit", SetOtherSmtpEdit)
	r.HandleFunc("/set/audit/list", SetAuditList)

	r.HandleFunc("/user/list", UserList)
	r.HandleFunc("/user/detail", UserDetail)
	r.HandleFunc("/user/set", UserSet)
	r.HandleFunc("/user/del", UserDel)
	r.HandleFunc("/user/online", UserOnline)
	r.HandleFunc("/user/offline", UserOffline)
	r.HandleFunc("/user/reline", UserReline)
	r.HandleFunc("/user/otp_qr", UserOtpQr)
	r.HandleFunc("/user/ip_map/list", UserIpMapList)
	r.HandleFunc("/user/ip_map/detail", UserIpMapDetail)
	r.HandleFunc("/user/ip_map/set", UserIpMapSet)
	r.HandleFunc("/user/ip_map/del", UserIpMapDel)

	r.HandleFunc("/group/list", GroupList)
	r.HandleFunc("/group/names", GroupNames)
	r.HandleFunc("/group/detail", GroupDetail)
	r.HandleFunc("/group/set", GroupSet)
	r.HandleFunc("/group/del", GroupDel)

	// pprof
	if base.Cfg.Pprof {
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline).Name("debug")
		r.HandleFunc("/debug/pprof/profile", pprof.Profile).Name("debug")
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol).Name("debug")
		r.HandleFunc("/debug/pprof/trace", pprof.Trace).Name("debug")
		r.HandleFunc("/debug/pprof", location("/debug/pprof/")).Name("debug")
		r.PathPrefix("/debug/pprof/").HandlerFunc(pprof.Index).Name("debug")
	}

	base.Info("Listen admin", base.Cfg.AdminAddr)

	// 修复 CVE-2016-2183
	cipherSuites := tls.CipherSuites()
	selectedCipherSuites := make([]uint16, 0, len(cipherSuites))
	for _, s := range cipherSuites {
		selectedCipherSuites = append(selectedCipherSuites, s.ID)
	}
	// 使用 https 协议
	// 设置tls信息
	tlsConfig := &tls.Config{
		NextProtos:   []string{"http/1.1"},
		MinVersion:   tls.VersionTLS12,
		CipherSuites: selectedCipherSuites,
	}
	srv := &http.Server{
		Addr:      base.Cfg.AdminAddr,
		Handler:   r,
		TLSConfig: tlsConfig,
	}
	err := srv.ListenAndServeTLS(base.Cfg.CertFile, base.Cfg.CertKey)
	if err != nil {
		base.Fatal(err)
	}
}

// 跳转
func location(url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
	}
}
