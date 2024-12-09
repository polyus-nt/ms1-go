package ms1

type BackTrackMsg struct {
	// TODO: CHECK uploadstage in user code
	UploadStage

	NoPacks bool

	CurPack    uint16
	TotalPacks uint16
}