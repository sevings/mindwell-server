package votes

import (
	"github.com/sevings/mindwell-server/restapi/operations/votes"
	"github.com/sevings/mindwell-server/utils"
)

// ConfigureAPI creates operations handlers
func ConfigureAPI(srv *utils.MindwellServer) {
	srv.API.VotesGetEntriesIDVoteHandler = votes.GetEntriesIDVoteHandlerFunc(newEntryVoteLoader(srv))
	srv.API.VotesPutEntriesIDVoteHandler = votes.PutEntriesIDVoteHandlerFunc(newEntryVoter(srv))
	srv.API.VotesDeleteEntriesIDVoteHandler = votes.DeleteEntriesIDVoteHandlerFunc(newEntryUnvoter(srv))

	srv.API.VotesGetCommentsIDVoteHandler = votes.GetCommentsIDVoteHandlerFunc(newCommentVoteLoader(srv))
	srv.API.VotesPutCommentsIDVoteHandler = votes.PutCommentsIDVoteHandlerFunc(newCommentVoter(srv))
	srv.API.VotesDeleteCommentsIDVoteHandler = votes.DeleteCommentsIDVoteHandlerFunc(newCommentUnvoter(srv))
}
