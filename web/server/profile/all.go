package profile

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/html-parser"
	"github.com/memocash/memo/app/profile"
	"github.com/memocash/memo/app/res"
	"net/http"
	"strings"
)

var allRoute = web.Route{
	Pattern: res.UrlProfiles,
	Handler: func(r *web.Response) {
		profilesByDate(r, true)
		r.RenderTemplate(res.TmplProfiles)
	},
}

func profilesByDate(r *web.Response, oldestToNewest bool) {
	r.Helper["Nav"] = "profiles"
	r.Helper["Title"] = "Memo - Profiles"
	offset := r.Request.GetUrlParameterInt("offset")
	searchString := html_parser.EscapeWithEmojis(r.Request.GetUrlParameter("s"))
	var selfPkHash []byte
	if auth.IsLoggedIn(r.Session.CookieId) {
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
			return
		}
		key, err := db.GetKeyForUser(user.Id)
		if err != nil {
			r.Error(jerr.Get("error getting key for user", err), http.StatusInternalServerError)
			return
		}
		selfPkHash = key.PkHash
	}
	var statOrderType db.UserStatOrderType
	if oldestToNewest {
		statOrderType = db.UserStatOrderCreated
	} else {
		statOrderType = db.UserStatOrderNewest
	}
	profiles, err := profile.GetProfiles(selfPkHash, searchString, offset, statOrderType)
	if err != nil {
		r.Error(jerr.Get("error getting profiles", err), http.StatusInternalServerError)
		return
	}
	err = profile.AttachReputationToProfiles(profiles)
	if err != nil {
		r.Error(jerr.Get("error attaching reputation to profiles", err), http.StatusInternalServerError)
		return
	}
	res.SetPageAndOffset(r, offset)
	r.Helper["SearchString"] = searchString
	var url string
	if oldestToNewest {
		url = res.UrlProfiles
	} else {
		url = res.UrlProfilesNew
	}
	if searchString != "" {
		r.Helper["OffsetLink"] = fmt.Sprintf("%s?s=%s", strings.TrimLeft(url, "/"), searchString)
	} else {
		r.Helper["OffsetLink"] = fmt.Sprintf("%s?", url)
	}
	r.Helper["Profiles"] = profiles
}
