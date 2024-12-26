package security

type TOTPService interface {
	GenerateCode() (string, error)
	ValidateCode(code string) (bool, error)
}
