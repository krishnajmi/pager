package login

const (
	UserTypeAdmin     = "admin"
	UserTypeMarketing = "marketing"
	UserTypeNormal    = "user"
)

const (
	PagerAdminAccess       = "PAGER.ADMIN"
	PagerNotifcationAccess = "PAGER.NOTIFICATION"
	PagerTemplateAccess    = "PAGER.CAMPAIGN_TRIGGER"
	PagerAuthAccess        = "PAGER.AUDIENCE"
)

var (
	DefaultAdminPermissions = []string{
		PagerAdminAccess,
	}

	DefaultMarketingPermissions = []string{
		PagerNotifcationAccess,
		PagerTemplateAccess,
	}

	DefaultUserPermissions = []string{
		PagerNotifcationAccess,
	}
)
