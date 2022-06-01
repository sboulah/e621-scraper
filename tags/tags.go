package tags

import "fmt"

func init() {
	fmt.Println("Father package initialized")
}

type RawPost struct {
	URL string `json:"url"`
}

type CleanedPost struct {
	POSTURL  string        `json:"post_url"`
	POSTID   int           `json:"id"`
	MEDIAURL string        `json:"media_url"`
	SOURCES  []interface{} `json:"sources"`
	MD5      string        `json:"md5"`
}

type MetaPost struct {
	ChangeSeq     int64         `json:"change_seq"`
	CommentCount  int64         `json:"comment_count"`
	CreatedAt     string        `json:"created_at"`
	Description   string        `json:"description"`
	FavCount      int64         `json:"fav_count"`
	File          File          `json:"file"`
	Flags         Flags         `json:"flags"`
	ID            int64         `json:"id"`
	IsFavorited   bool          `json:"is_favorited"`
	LockedTags    []interface{} `json:"locked_tags"`
	Pools         []interface{} `json:"pools"`
	Relationships Relationships `json:"relationships"`
	Sample        Sample        `json:"sample"`
	Score         Score         `json:"score"`
	Sources       []interface{} `json:"sources"`
	Tags          []interface{} `json:"tags"`
	UpdatedAt     string        `json:"updated_at"`
	UploaderID    int64         `json:"uploader_id"`
}

type File struct {
	EXT    string `json:"ext"`
	Height int64  `json:"height"`
	Md5    string `json:"md5"`
	Size   int64  `json:"size"`
	URL    string `json:"url"`
	Width  int64  `json:"width"`
}

type Flags struct {
	Deleted      bool `json:"deleted"`
	Flagged      bool `json:"flagged"`
	HasNotes     bool `json:"has_notes"`
	NoteLocked   bool `json:"note_locked"`
	Pending      bool `json:"pending"`
	RatingLocked bool `json:"rating_locked"`
	StatusLocked bool `json:"status_locked"`
}

type Relationships struct {
	Children          []interface{} `json:"children"`
	HasActiveChildren bool          `json:"has_active_children"`
	HasChildren       bool          `json:"has_children"`
	ParentID          interface{}   `json:"parent_id"`
}

type Sample struct {
	Alternates Alternates `json:"alternates"`
	Has        bool       `json:"has"`
	Height     int64      `json:"height"`
	URL        string     `json:"url"`
	Width      int64      `json:"width"`
}

type Alternates struct {
}

type Score struct {
	Down  int64 `json:"down"`
	Total int64 `json:"total"`
	Up    int64 `json:"up"`
}
