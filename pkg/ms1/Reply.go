package ms1

import (
	"bytes"
	"fmt"
	"github.com/polyus-nt/ms1-go/internal/xxd"
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
	Nanoid string
}

func (id ID) String() string {
	return fmt.Sprintf("ID [ %#x ] : [ 0x%v ]", id.Mark, id.Nanoid)
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

type Meta struct {
	Mark          int
	Valid         bool
	RefBlHw       string // Описывает физическое окружение контроллера (плату)
	RefBlFw       string // Указывает на версию прошивки загрузчика
	RefBlUserCode string //
	RefBlChip     string // Указывает на контроллер, здесь то, что нужно для компиляции прошивки
	RefBlProtocol string // Описывает возможности протокола загрузчика
	RefCgHw       string // Указывает на аппаратное исполнение
	RefCgFw       string // Указывает на версию прошивки кибергена
	RefCgProtocol string // Указывает на возможности протокола кибергена
}

func (m Meta) String() string {

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("Meta [ %#x ] -> {\n\tValid: %v\n\trefBlHw: %v\n\trefBlFw : %v\n\trefBlUserCode: %v\n\trefBlChip : %v\n\trefBlProtocol: %v\n\trefCgHw: %v\n\trefCgFw: %v\n\trefCgProtocol: %v\n}", m.Mark, m.Valid, m.RefBlHw, m.RefBlFw, m.RefBlUserCode, m.RefBlChip, m.RefBlProtocol, m.RefCgHw, m.RefCgFw, m.RefCgProtocol))

	return buf.String()
}

func (m Meta) GetStatus() ReplyStatus {
	return REPLY_META
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
