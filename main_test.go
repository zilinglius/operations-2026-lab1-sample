package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthz(t *testing.T) {
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest("GET", "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("/healthz 返回 %d，期望 200", rec.Code)
	}
}

func TestVersionEndpoint(t *testing.T) {
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest("GET", "/version", nil))
	var got versionInfo
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("/version 不是合法 JSON：%v", err)
	}
	if got.Version == "" || got.Commit == "" || got.BuildNo == "" {
		t.Errorf("/version 缺少字段：%+v", got)
	}
}

func TestConflictEndpoint(t *testing.T) {
	post := func(body string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/conflict", strings.NewReader(body))
		newMux().ServeHTTP(rec, req)
		return rec
	}

	// 单周与双周在同一时段上课，不冲突。
	rec := post(`{"a":{"id":"A","weekday":1,"start":1,"end":2,"weeks":"odd"},"b":{"id":"B","weekday":1,"start":1,"end":2,"weeks":"even"}}`)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"conflict":false`) {
		t.Errorf("单双周同一时段：状态 %d，响应 %q", rec.Code, rec.Body.String())
	}

	// 每周上课的两个班同一时段，冲突。
	rec = post(`{"a":{"id":"A","weekday":1,"start":1,"end":2,"weeks":"all"},"b":{"id":"B","weekday":1,"start":2,"end":3,"weeks":"all"}}`)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"conflict":true`) {
		t.Errorf("节次重叠：状态 %d，响应 %q", rec.Code, rec.Body.String())
	}

	if rec := post(`{"a":`); rec.Code != http.StatusBadRequest {
		t.Errorf("坏 JSON 应返回 400，实际 %d", rec.Code)
	}
	if rec := post(`{"a":{"weekday":9,"start":1,"end":2,"weeks":"all"},"b":{"weekday":1,"start":1,"end":2,"weeks":"all"}}`); rec.Code != http.StatusBadRequest {
		t.Errorf("不合法时段应返回 400，实际 %d", rec.Code)
	}
}
