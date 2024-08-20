package internal

// general type (for enum impl)
type PacketWrite interface{}

// derrived types for PacketWrite
type PwFrame struct {
	frame Frame
}

type PwNuke struct {
	index uint8
}

type PwStatusApp struct {
	fwId   string
	gitRev string
}

type PwReset struct{}

// general type
type PacketRead interface{}

// derrived types for PacketRead
type PrAck struct{}

type PrError struct {
	value string
}

type PrDone struct{}

type PrFrame struct {
	frame Frame
}

type PrStatusApp struct {
	fwId   string
	gitRev string
}

type PrStatusOwn struct {
	fwId   string
	gitRev string
	e1     string
}