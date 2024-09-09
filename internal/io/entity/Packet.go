package entity

type Packet struct {
	Mark uint8
	Addr Address
	Code string
	Load []Load
}