package role

type Role string

const (
	Undefined Role = "\n"
	User      Role = "User"
	Moderator Role = "Moderator"
	Admin     Role = "Admin"
)
