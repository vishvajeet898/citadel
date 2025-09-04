package constants

const (
	USER_TYPE_SUPER_ADMIN    = "super-admin"
	USER_TYPE_SYSTEM         = "system"
	USER_TYPE_ADMIN          = "admin"
	USER_TYPE_PATHOLOGIST    = "pathologist"
	USER_TYPE_LAB_TECHNICIAN = "lab-technician"
	USER_TYPE_CRM            = "crm"
	USER_TYPE_ACCESSIONER    = "accessioner"
)

var USER_TYPES = []string{
	USER_TYPE_SUPER_ADMIN,
	USER_TYPE_SYSTEM,
	USER_TYPE_ADMIN,
	USER_TYPE_PATHOLOGIST,
	USER_TYPE_LAB_TECHNICIAN,
	USER_TYPE_CRM,
	USER_TYPE_ACCESSIONER,
}
