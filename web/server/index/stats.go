package index

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var statsRoute = web.Route{
	Pattern: res.UrlStats,
	Handler: func(r *web.Response) {
		memoFollowCount, err := db.GetCountMemoFollows()
		if err != nil {
			r.Error(jerr.Get("error getting memo follow count", err), http.StatusInternalServerError)
			return
		}
		memoLikeCount, err := db.GetCountMemoLikes()
		if err != nil {
			r.Error(jerr.Get("error getting memo like count", err), http.StatusInternalServerError)
			return
		}
		memoPostCount, memoVotePostCount, memoTopicPostCount, memoReplyPostCount, err := db.GetCountMemoPosts()
		if err != nil {
			r.Error(jerr.Get("error getting memo post count", err), http.StatusInternalServerError)
			return
		}
		memoSetNameCount, err := db.GetCountMemoSetName()
		if err != nil {
			r.Error(jerr.Get("error getting memo set name count", err), http.StatusInternalServerError)
			return
		}
		memoSetProfileCount, err := db.GetCountMemoSetProfile()
		if err != nil {
			r.Error(jerr.Get("error getting memo set profile count", err), http.StatusInternalServerError)
			return
		}
		memoSetProfilePicCount, err := db.GetCountMemoSetPic()
		if err != nil {
			r.Error(jerr.Get("error getting memo set pic count", err), http.StatusInternalServerError)
			return
		}
		memoPollQuestionCount, err := db.GetCountMemoPollQuestion()
		if err != nil {
			r.Error(jerr.Get("error getting memo poll question count", err), http.StatusInternalServerError)
			return
		}
		memoPollOptionCount, err := db.GetCountMemoPollOption()
		if err != nil {
			r.Error(jerr.Get("error getting memo poll option count", err), http.StatusInternalServerError)
			return
		}
		memoPollVoteCount, err := db.GetCountMemoPollVote()
		if err != nil {
			r.Error(jerr.Get("error getting memo poll vote count", err), http.StatusInternalServerError)
			return
		}
		memoTopicFollowCount, err := db.GetCountMemoTopicFollow()
		if err != nil {
			r.Error(jerr.Get("error getting memo topic follow count", err), http.StatusInternalServerError)
			return
		}
		uniqueUsers, err := db.GetUniqueUserCount()
		if err != nil {
			r.Error(jerr.Get("error getting unique users", err), http.StatusInternalServerError)
			return
		}
		r.Helper["MemoFollowCount"] = int64(memoFollowCount)
		r.Helper["MemoLikeCount"] = int64(memoLikeCount)
		r.Helper["MemoPostCount"] = int64(memoPostCount)
		r.Helper["MemoVotePostCount"] = int64(memoVotePostCount)
		r.Helper["MemoReplyPostCount"] = int64(memoReplyPostCount)
		r.Helper["MemoTopicPostCount"] = int64(memoTopicPostCount)
		r.Helper["MemoSetNameCount"] = int64(memoSetNameCount)
		r.Helper["MemoSetProfileCount"] = int64(memoSetProfileCount)
		r.Helper["MemoSetProfilePicCount"] = int64(memoSetProfilePicCount)
		r.Helper["MemoPollQuestionCount"] = int64(memoPollQuestionCount)
		r.Helper["MemoPollOptionCount"] = int64(memoPollOptionCount)
		r.Helper["MemoPollVoteCount"] = int64(memoPollVoteCount)
		r.Helper["MemoTopicFollowCount"] = int64(memoTopicFollowCount)
		r.Helper["MemoTotalPosts"] = int64(memoPostCount) +
			int64(memoPollQuestionCount) +
			int64(memoReplyPostCount) +
			int64(memoVotePostCount) +
			int64(memoTopicPostCount)
		r.Helper["MemoTotalActionCount"] = int64(memoFollowCount +
			memoLikeCount +
			memoPostCount +
			memoReplyPostCount +
			memoTopicPostCount +
			memoSetNameCount +
			memoSetProfileCount +
			memoSetProfilePicCount +
			memoPollQuestionCount +
			memoPollOptionCount +
			memoPollVoteCount +
			memoTopicFollowCount)
		r.Helper["UniqueUsers"] = uniqueUsers

		r.Helper["Title"] = "Memo - Stats"

		r.RenderTemplate(res.TmplStats)
	},
}
