package repository

type Repository interface {
	saveMsg(msg string) error
	getAllMsg() (string, error)
}
