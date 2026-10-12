// enroll-conflict 是实验 1 使用的最小示例服务：判断两个教学班的上课时间是否冲突。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

// 以下三个变量在构建时由 -ldflags "-X main.xxx=..." 注入（见 Makefile）。
// 没有注入时保持默认值，这也是“本地随手 go build”得到的二进制的样子。
var (
	version = "dev"
	commit  = "unknown"
	buildNo = "local"
)

type versionInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	BuildNo string `json:"build_no"`
	Go      string `json:"go"`
}

func currentVersion() versionInfo {
	return versionInfo{Version: version, Commit: commit, BuildNo: buildNo, Go: runtime.Version()}
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()

	// 存活检查：进程能响应即可（L-05 第 23 页）。
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	// “线上跑的是什么”：一次请求回答版本号、提交号、构建号。
	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(currentVersion())
	})

	mux.HandleFunc("POST /api/conflict", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			A Section `json:"a"`
			B Section `json:"b"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "请求体不是合法的 JSON", http.StatusBadRequest)
			return
		}
		if !req.A.Valid() || !req.B.Valid() {
			http.Error(w, "时段不合法", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"conflict": Conflicts(req.A, req.B)})
	})

	return mux
}

func main() {
	addr := flag.String("addr", ":8080", "监听地址")
	showVersion := flag.Bool("version", false, "打印版本信息后退出")
	flag.Parse()

	if *showVersion {
		v := currentVersion()
		fmt.Printf("enroll-conflict version=%s commit=%s build_no=%s go=%s\n", v.Version, v.Commit, v.BuildNo, v.Go)
		return
	}

	log.Printf("enroll-conflict %s (%s) 监听 %s", version, commit, *addr)
	log.Fatal(http.ListenAndServe(*addr, newMux()))
}
