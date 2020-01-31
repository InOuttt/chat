package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/abaron/chat/server/store"
	"github.com/abaron/chat/server/store/types"
)

func handleMessageRemove(w http.ResponseWriter, r *http.Request) {
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Basic" {
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 || !validateBasicAuth(pair[0], pair[1]) {
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	type Response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}

	var (
		from int64
		to   int64
		resp Response
	)

	keys, ok := r.URL.Query()["from"]
	if !ok || len(keys[0]) < 1 {
		from = 0
	} else {
		i, err := strconv.ParseInt(keys[0], 10, 64)
		if err != nil {
			from = 0
		} else {
			from = i
		}
	}

	keys, ok = r.URL.Query()["to"]
	if !ok || len(keys[0]) < 1 {
		to = 0
	} else {
		i, err := strconv.ParseInt(keys[0], 10, 64)
		if err != nil {
			to = 0
		} else {
			to = i
		}
	}

	err := store.MessageRemoveAPI(from, to)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp = Response{
		Status:  true,
		Message: "Success",
	}
	json.NewEncoder(w).Encode(resp)
}

func handleMessageUnreadCount(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	if len(uid) == 0 {
		http.Error(
			w,
			http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized,
		)
		return
	}

	uuid, _ := strconv.ParseUint(uid, 10, 64)
	if uuid == 0 {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := struct {
		Status  bool
		Message string
		Unread  int
	}{}

	var err error
	resp.Unread, err = store.Users.GetUnreadCount(types.Uid(uuid))
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}

	fmt.Fprint(w, resp.Unread)
}

func validateBasicAuth(username, password string) bool {
	if username == "baron" && password == "baron1234567890-=!@#$%^&*()_+" {
		return true
	}
	return false
}
