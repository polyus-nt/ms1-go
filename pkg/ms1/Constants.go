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
	REPLY_META
	REPLY_GARBAGE
	REPLY_ERROR
)
