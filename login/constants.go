package login

const (
	UserTypeAdmin     = "admin"
	UserTypeMarketing = "marketing"
	UserTypeNormal    = "user"
)

const (
	PagerAdminAccess           = "PAGER.ADMIN"
	PagerCampaignAccess        = "PAGER.CAMPAIGN"
	PagerCampaignTriggerAccess = "PAGER.CAMPAIGN_TRIGGER"
	PagerAudienceAccess        = "PAGER.AUDIENCE"
	PagerAudienceEditAccess    = "PAGER.AUDIENCE_EDIT"
	PagerWebhookAccess         = "PAGER.WEBHOOK"
	PagerCommunicationAccess   = "PAGER.COMMUNICATION"
	PagerDNAccess              = "PAGER.DN"
	PagerGroupDNAccess         = "PAGER.GROUP_DN"
	PagerInternalAccess        = "PAGER.INTERNAL"
	PagerJourneyAccess         = "PAGER.JOURNEY"
	PagerReportAccess          = "PAGER.REPORT"
	PagerTemplateAccess        = "PAGER.TEMPLATE"
	PagerTemplateEditAccess    = "PAGER.TEMPLATE_EDIT"
	PagerTenantAccess          = "PAGER.TENANT"
	PagerUserAccess            = "PAGER.USER"
)

var (
	DefaultAdminPermissions = []string{
		PagerAdminAccess,
	}

	DefaultMarketingPermissions = []string{
		PagerCampaignAccess,
		PagerCampaignTriggerAccess,
		PagerAudienceAccess,
		PagerAudienceEditAccess,
		PagerTemplateAccess,
		PagerReportAccess,
		PagerCommunicationAccess,
		PagerTemplateEditAccess,
	}

	DefaultUserPermissions = []string{
		PagerCampaignAccess,
		PagerAudienceAccess,
		PagerReportAccess,
	}
)
