package errors

import "errors"

var (
	ErrListenAndServeServer = errors.New("serving server")
	ErrShutdownServer       = errors.New("shutdown server")

	ErrConnectDB   = errors.New("connect to database")
	ErrExecQuery   = errors.New("executing query")
	ErrTransaction = errors.New("executing transaction")
	ErrCommit      = errors.New("commit error")

	ErrMigrationFailed = errors.New("migration failed")
	ErrValidateFailed  = errors.New("validation failed")

	ErrDuplicateWarehouse = errors.New("duplicate warehouse")
	ErrDuplicateRoute     = errors.New("duplicate route")
	ErrDuplicateUser      = errors.New("duplicate user")
	ErrDuplicateThreshold = errors.New("duplicate threshold")

	ErrNotFoundWarehouse  = errors.New("warehouse not found")
	ErrNotFoundRoute      = errors.New("route not found")
	ErrNotFoundUser       = errors.New("user not found")
	ErrNotFoundDriver     = errors.New("driver not found")
	ErrNotFoundThreshold  = errors.New("threshold not found")
	ErrNotFoundTruckCall  = errors.New("truck call not found")

	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrForbidden          = errors.New("forbidden")

	ErrMLServiceUnavailable = errors.New("ml service unavailable")
	ErrMLServiceError       = errors.New("ml service error")
)