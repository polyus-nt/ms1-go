package ms1

type BackTrackMsg struct {
	UploadStage

	NoPacks bool

	CurPack    uint16
	TotalPacks uint16
}