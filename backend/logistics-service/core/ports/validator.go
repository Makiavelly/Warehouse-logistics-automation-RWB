package ports

type Validator interface {
	ValidateStruct(s any) error
}