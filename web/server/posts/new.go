package posts

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

var newRoute = web.Route{
	Pattern: res.UrlPostsNew,
	Handler: func(r *web.Response) {
		preHandler(r)
		offset := r.Request.GetUrlParameterInt("offset")
		searchString := html_parser.EscapeWithEmojis(r.Request.GetUrlParameter("s"))
		var userPkHash []byte
		var userId uint
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
			userPkHash = key.PkHash
			userId = user.Id
		}
		posts, err := profile.GetRecentPosts(userPkHash, uint(offset), searchString)
		if err != nil {
			r.Error(jerr.Get("error getting recent posts", err), http.StatusInternalServerError)
			return
		}
		err = profile.AttachParentToPosts(posts)
		if err != nil {
			r.Error(jerr.Get("error attaching parent to posts", err), http.StatusInternalServerError)
			return
		}
		err = profile.AttachLikesToPosts(posts)
		if err != nil {
			r.Error(jerr.Get("error attaching likes to posts", err), http.StatusInternalServerError)
			return
		}
		err = profile.AttachPollsToPosts(posts)
		if err != nil {
			r.Error(jerr.Get("error attaching polls to posts", err), http.StatusInternalServerError)
			return
		}
		if len(userPkHash) > 0 {
			err = profile.AttachReputationToPosts(posts)
			if err != nil {
				r.Error(jerr.Get("error attaching reputation to posts", err), http.StatusInternalServerError)
				return
			}
		}
		err = profile.SetShowMediaForPosts(posts, userId)
		if err != nil {
			r.Error(jerr.Get("error setting show media for posts", err), http.StatusInternalServerError)
			return
		}
		res.SetPageAndOffset(r, offset)
		r.Helper["Posts"] = posts
		r.Helper["Title"] = "Memo - New Posts"
		r.Helper["SearchString"] = searchString
		if searchString != "" {
			r.Helper["OffsetLink"] = fmt.Sprintf("%s?s=%s", strings.TrimLeft(res.UrlPostsNew, "/"), searchString)
		} else {
			r.Helper["OffsetLink"] = fmt.Sprintf("%s?", res.UrlPostsNew)
		}
		r.Render()
	},
}
