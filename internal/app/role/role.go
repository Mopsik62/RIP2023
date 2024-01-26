package role

type Role string

const (
	Undefined Role = "Undefined"
	User      Role = "User"
	Moderator Role = "Moderator"
	Admin     Role = "Admin"
)
