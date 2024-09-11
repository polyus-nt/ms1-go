package ms1

import (
	"bytes"
	"fmt"
	"ms1-go/internal/xxd"
)

// Reply - тип интерфейс для проброса ответа пользователю
type Reply interface {
	fmt.Stringer
	GetStatus() ReplyStatus
}

type Ping struct{ Value int }

func (p Ping) String() string {
	return fmt.Sprintf("Ping [%#x]", p.Value)
}
func (p Ping) GetStatus() ReplyStatus {
	return REPLY_PING
}

type Pong struct{ Value int }

func (p Pong) String() string {
	return fmt.Sprintf("Pong [%#x]", p.Value)
}
func (p Pong) GetStatus() ReplyStatus {
	return REPLY_PONG
}

type GenePong struct{ Value int }

func (gp GenePong) String() string {
	return fmt.Sprintf("Pong from gene [%#x]", gp.Value)
}
func (gp GenePong) GetStatus() ReplyStatus {
	return REPLY_GENE_PONG
}

type GeneAck struct{ Value int }

func (ga GeneAck) String() string {
	return fmt.Sprintf("Ack from gene [%#x]", ga.Value)
}
func (ga GeneAck) GetStatus() ReplyStatus {
	return REPLY_GENE_ACK
}

type Ack struct{ Value int }

func (a Ack) String() string {
	return fmt.Sprintf("Ack [%#x]", a.Value)
}
func (a Ack) GetStatus() ReplyStatus {
	return REPLY_ACK
}

type Nack struct{ Value int }

func (n Nack) String() string {
	return fmt.Sprintf("Nack [%#x]", n.Value)
}
func (n Nack) GetStatus() ReplyStatus {
	return REPLY_NACK
}

type Ref struct{ Value int64 }

func (r Ref) String() string {
	return fmt.Sprintf("Ref [%#x]", r.Value)
}
func (r Ref) GetStatus() ReplyStatus {
	return REPLY_REF
}

type ID struct {
	Mark   int
	Nanoid int64
}

func (id ID) String() string {
	return fmt.Sprintf("ID [ %#x ] : [ %#x ]", id.Mark, id.Nanoid)
}
func (id ID) GetStatus() ReplyStatus {
	return REPLY_ID
}

type Frame2 struct {
	Page  int
	Index int
	Mark  int
	Blob  string
}

func (f Frame2) String() string {

	buf := bytes.Buffer{}

	buf.WriteString(fmt.Sprintf("Frame [ %#x ] %v.%v\n", f.Mark, f.Page, f.Index))

	var bin []xxd.Bin
	remains := []byte(f.Blob)
	for len(remains) > 0 {
		I := max(len(remains)-16, 0)
		qword := remains[I:]
		var res [][]byte
		for i := len(qword) - 2; i >= 0; i -= 2 {
			res = append(res, qword[i:i+2])
		}
		bin = append(bin, xxd.Bin(bytes.Join(res, []byte(""))))
		remains = remains[:I]
	}
	buf.WriteString(xxd.Xxd(bin))

	return buf.String()
}
func (f Frame2) GetStatus() ReplyStatus {
	return REPLY_FRAME
}

type Garbage struct {
	Comment string
	Garbage string
}

func (g Garbage) String() string {
	return fmt.Sprintf("Garbage %v : %v", g.Comment, g.Garbage)
}
func (g Garbage) GetStatus() ReplyStatus {
	return REPLY_GARBAGE
}

type Error struct {
	Mark    int
	Message string
}

func (e Error) String() string {
	return fmt.Sprintf("Error %v : %v", e.Mark, e.Message)
}
func (e Error) GetStatus() ReplyStatus {
	return REPLY_ERROR
}