package config

const (
	CfgFileFlag     = "cfg"
	CfgStdinFlag    = "xargs"
	DefaultFlag     = "default"
	CronCmd         = "cron"
	RemoveCmd       = "remove"
	InstallCmd      = "install"
	VersionCmd      = "version"
	Health          = "health"
	AllFlag         = "all"
	OptionsFlag     = "options"
	LogRedirectFlag = "log-redirect"
	////////////////////////////////// article commands
	NewsCmd           = "articles"
	DisplayOptionsCmd = "options"
	DisplayIdsCmd     = "ids"
	FetchFlagCmd      = "fetch"
	PortalIdFlag      = "id"
	SettingsLoaded    = "SettingsLoaded"
)

// //////////////////////////////// Fetch const
const (
	ActionOutSubDir       = "out"
	ActionPage            = "page"
	ActionLimitBatchSize  = "batch-size"
	ActionLimitSortByDate = "date"
	ActionSectionFlag     = "section"
	ActionQueryValue      = "query"
	ActionEndIndex        = "endIndex"
)

// paths
const (
	ConfigurationFile = "CONFIGURATION_FILE"
	InstallationPath  = "{{__INSTALLATION_PLACEHOLDER__}}"
	//InstallationPath  = "/usr/local/bin"
	Command = "{{__COMMAND_PLACEHOLDER__}}"
	Company = "{{__COMPANY_PLACEHOLDER__}}"
)

// Config paths
var (
	CronScriptPrefix = "crontab_%v.sh"
)
