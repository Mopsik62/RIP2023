package ds

type Substances struct {
	ID     uint `gorm:"primarykey"`
	Title  string
	Text   string
	Number int
	Image  string `gorm:"type:bytea"`
	Status string
}
