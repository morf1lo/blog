package service

type ActivationMailData struct {
	To   []string
	Link string
}

type ResetPasswordTokenMailData struct {
	To    []string
	Token string
}
