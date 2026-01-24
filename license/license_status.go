package license

type LicenseStatus string

const (
	OK        LicenseStatus = "OK"
	SOFT_LOCK LicenseStatus = "SOFT_LOCK"
	HARD_LOCK LicenseStatus = "HARD_LOCK"
)
