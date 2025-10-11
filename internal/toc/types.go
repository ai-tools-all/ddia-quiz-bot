package toc

type Header struct {
	Level   int
	Title   string
	Anchor  string
	LineNum int
}

type Options struct {
	MaxDepth          int
	Format            string
	Recursive         bool
	OutputFile        string
	ConfirmThreshold  int
	SkipConfirmation  bool
}
