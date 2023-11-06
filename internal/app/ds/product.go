package ds

type Substances struct {
	ID      uint `gorm:"primarykey"`
	Title   string
	Class   string
	Formula string
	Image   string `gorm:"type:bytea"`
	Status  string
}
