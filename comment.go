package esprimago

type CommentType int

const (
	Block CommentType = iota
	Line
)

type Comment struct {
	Type      CommentType
	Value     *string
	MultiLine bool
	Slice     []int
	Start     int
	End       int
	Loc       *SourceLocation
}
