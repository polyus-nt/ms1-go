package ms1

type ReplyStatus = uint8

const ZeroAddress = "0000000000000000"

//goland:noinspection ALL
const (
	REPLY_PING ReplyStatus = iota
	REPLY_PONG
	REPLY_GENE_PONG
	REPLY_GENE_ACK
	REPLY_ACK
	REPLY_NACK
	REPLY_REF
	REPLY_ID
	REPLY_FRAME
	REPLY_GARBAGE
	REPLY_ERROR
)

// TODO: add map with device info by REF_ID (MDI)
// TODO: add iota generator const for nanoid about devices (and use later as key for map with ms1 info (MDI))