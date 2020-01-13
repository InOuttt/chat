package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/abaron/chat/server/auth"
	"github.com/abaron/chat/server/store"
	"github.com/abaron/chat/server/store/types"
)

func handleLdapExchangeToken(w http.ResponseWriter, r *http.Request) {
	const contentType = "application/json"
	var (
		user    types.User
		ldapURL = globals.ldapServer["base_url"] + "/api/ldap/chat"
	)
	ldapToken := r.URL.Query().Get("token")
	if len(ldapToken) == 0 {
		http.Error(
			w,
			http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized,
		)
		return
	}

	// step 1. check if user already registered
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(map[string]string{
		"action":        "ldap_check",
		"client_id":     globals.ldapServer["client_id"],
		"client_secret": globals.ldapServer["client_secret"],
		"token":         ldapToken,
	})

	resp, err := http.Post(ldapURL, contentType, &buf)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(
			w,
			http.StatusText(resp.StatusCode),
			resp.StatusCode,
		)
		return
	}

	m := struct {
		ID     uint64  `json:"id"`
		Name   string  `json:"name"`
		LdapID *string `json:"ldap_id"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	if m.LdapID != nil && *m.LdapID != "" {
		ldapID, _ := strconv.ParseUint(*m.LdapID, 10, 64)
		if ldapID == 0 {
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}

		user.SetUid(types.Uid(ldapID))
		goto here
	}

	// step 2. if not registered, then register new user & link it
	// step 2a. register new user
	user.Access = types.DefaultAccess{Auth: types.ModeCAuth, Anon: types.ModeNone}
	user.Public = map[string]string{"fn": m.Name}
	_, err = store.Users.Create(&user, nil)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	// step 2b. link
	buf.Reset()
	json.NewEncoder(&buf).Encode(map[string]interface{}{
		"action":        "ldap_link",
		"client_id":     globals.ldapServer["client_id"],
		"client_secret": globals.ldapServer["client_secret"],
		"token":         ldapToken,
		"ldap_id":       strconv.FormatUint(uint64(user.Uid()), 10),
	})

	resp, err = http.Post(ldapURL, contentType, &buf)
	if err != nil {
		store.Users.Delete(user.Uid(), true)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		store.Users.Delete(user.Uid(), true)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

here:
	// step 3. generate user token
	token, time, _ := store.GetLogicalAuthHandler("token").GenSecret(&auth.Rec{
		Uid:       user.Uid(),
		AuthLevel: auth.LevelAuth,
		Lifetime:  time.Hour * 24,
		Features:  auth.FeatureValidated})

	encToken, _ := json.Marshal(map[string]interface{}{
		"token":   token,
		"expires": time,
	})

	var qs = make(url.Values, 1)
	qs.Add("token", base64.RawStdEncoding.EncodeToString(encToken))
	http.Redirect(w, r, fmt.Sprintf("/#?%s", qs.Encode()), http.StatusSeeOther)
}
