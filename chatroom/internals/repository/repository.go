package repository

type Repository interface {
	SaveMsg(string) error
	GetAllMsg() ([]string, error)
}
