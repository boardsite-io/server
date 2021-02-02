package types

type BoardRequest struct {
	Action string `json:"action"`
}

type CreateBoardResponse struct {
	ID string `json:"id"`
}

type PageRequestData struct {
	PageID string `json:"pageId"`
	Index  int    `json:"index"`
}

type PageRankResponse struct {
	PageRank []string `json:"pageRank"`
}
