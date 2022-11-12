package repository

type Repository interface {
	SaveMsg(string, string) error
	GetAllMsg() (string, error)
}
