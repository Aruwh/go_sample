package common

type (
	// LogLevel represents different log levels
	LogLevel string
	// RequestPermission represents different request permissions
	RequestPermission string
	// AdminUserType represents different administrator types
	AdminUserType int
	// OrderType represents different variants of order types
	OrderType string
	// SortType represents different variants of sortBy types
	SortByType string
	// SaisonType represents different variants of the on going saisons
	SaisonType int
	// BookingtStatus represents different variants of the booking status
	BookingtStatus int
	// AudianceType represents different variants of audiances used in the JWT token
	AudianceType string
	// Locale represents different varaints of locales which used by users and admin users
	Locale string
	// Action represents different varaints of action which used to create a process_log entry
	Action string
	// Domain represents different varaints of the domains
	Domain string
	// Salutation represents different varaints of a title
	Sex int
	// PhoneNumberType represents different varaints of the a phone number
	PhoneNumberType int
	// PictureVariant replresents the different varaints of a pricture (file)
	PictureVariant string
)

const (
	SUPER_ADMINISTRATOR AdminUserType = iota
	ADMINISTRATOR
	APARTMENT_OWNER
)

const (
	SaisonPeak SaisonType = iota
	SaisonHigh
	SaisonMiddle
	SaisonLow
	SaisonBase
)

const (
	Mobile PhoneNumberType = iota
	Standard
)

const (
	Male Sex = iota
	Female
	Diverse
)

const (
	Available BookingtStatus = iota
	Reserved
	Confirmed
	Canceled
	Completed
	Blocked
	BlockedByAdmin
)

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	PANIC LogLevel = "PANIC"

	ADMIN_USER_VIEW    RequestPermission = "ADMIN_USER_VIEW"
	ADMIN_USER_EDIT    RequestPermission = "ADMIN_USER_EDIT"
	ADMIN_USER_DELETE  RequestPermission = "ADMIN_USER_DELETE"
	APARTMENT_VIEW     RequestPermission = "APARTMENT_VIEW"
	APARTMENT_EDIT     RequestPermission = "APARTMENT_EDIT"
	APARTMENT_DELETE   RequestPermission = "APARTMENT_DELETE"
	BOOKING_VIEW       RequestPermission = "BOOKING_VIEW"
	BOOKING_EDIT       RequestPermission = "BOOKING_EDIT"
	BOOKING_DELETE     RequestPermission = "BOOKING_DELETE"
	PERMISSION_VIEW    RequestPermission = "PERMISSION_VIEW"
	PERMISSION_EDIT    RequestPermission = "PERMISSION_EDIT"
	PERMISSION_DELETE  RequestPermission = "PERMISSION_DELETE"
	REAL_ESTATE_VIEW   RequestPermission = "REAL_ESTATE_VIEW"
	REAL_ESTATE_EDIT   RequestPermission = "REAL_ESTATE_EDIT"
	REAL_ESTATE_DELETE RequestPermission = "REAL_ESTATE_DELETE"
	NEWS_VIEW          RequestPermission = "NEWS_VIEW"
	NEWS_EDIT          RequestPermission = "NEWS_EDIT"
	NEWS_DELETE        RequestPermission = "NEWS_DELETE"
	SETTINGS_VIEW      RequestPermission = "SETTINGS_VIEW"
	SETTINGS_EDIT      RequestPermission = "SETTINGS_EDIT"
	SETTINGS_DELETE    RequestPermission = "SETTINGS_DELETE"
	USER_VIEW          RequestPermission = "USER_VIEW"
	USER_EDIT          RequestPermission = "USER_EDIT"
	USER_DELETE        RequestPermission = "USER_DELETE"
	PICTURE_VIEW       RequestPermission = "PICTURE_VIEW"
	PICTURE_EDIT       RequestPermission = "PICTURE_EDIT"
	PICTURE_DELETE     RequestPermission = "PICTURE_DELETE"
	TEST_VIEW          RequestPermission = "TEST_VIEW"
	TEST_EDIT          RequestPermission = "TEST_EDIT"
	TEST_DELETE        RequestPermission = "TEST_DELETE"

	AUDIANCE_USER             AudianceType = "U"
	AUDIANCE_APARTMENT_OWNER  AudianceType = "O"
	AUDIANCE_ADMIN_USER       AudianceType = "A"
	AUDIANCE_SUPER_ADMIN_USER AudianceType = "SA"

	SuperAdminEmail = "fewoserv@admin.ch"

	IdentityIdentifier    string = "IdentityIdentifier"
	CorrelationIdentifier string = "CorrelationIdentifier"

	OrderASC  OrderType = "asc"
	OrderDESC OrderType = "desc"

	SortByID      SortByType = "_id"
	SortByName    SortByType = "name"
	SortByCreated SortByType = "created.time"

	DeDE Locale = "deDE"
	EnGB Locale = "enGB"
	FrFR Locale = "frFR"
	ItIT Locale = "itIT"

	CREATED  Action = "CREATED"
	UPDATED  Action = "UPDATED"
	DELETED  Action = "DELETED"
	SEND     Action = "SEND"
	ACCEPTED Action = "ACCEPTED"

	ADMIN_USER  Domain = "ADMIN_USER"
	APARTMENT   Domain = "APARTMENT"
	BOOKING     Domain = "BOOKING"
	PERMISSION  Domain = "PERMISSION"
	REAL_ESTATE Domain = "REAL_ESTATE"
	PICTURE     Domain = "PICTURE"
	ATTRIBUTE   Domain = "ATTRIBUTE"
	SEASON      Domain = "SEASON"
	NEWS        Domain = "NEWS"
	USER        Domain = "USER"

	SMALL  PictureVariant = "small"
	MIDDLE PictureVariant = "middle"
	LARGE  PictureVariant = "large"
	ORIGIN PictureVariant = "origin"
)
