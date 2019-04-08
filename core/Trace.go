package core

// TraceChange is a boolean that, when true, turns on the tracing of ChangeNotifications
var TraceChange bool

// OmitHousekeepingCalls is a boolean that indicates whether housekeeping calls should be omitted when TraceChanve is enabled.
var OmitHousekeepingCalls bool

// OmitManageTreeNodesCalls is a boolean that indicates whether tree node management calls should be omitted when TraceChanve is enabled.
var OmitManageTreeNodesCalls bool

// TraceLocks is a boolean that, when true, turns on the tracing of locks
var TraceLocks = false

// TraceNotifications determines whether individual notifications will be traced. Its primary purpose
// is to track down cycles in notifications. When true, every time a notification is created a graph of
// the notification and its antecedents will be created and added to the notificationGraphs. In addition,
// if enableNotificationPrint is true, every time a notification is created it and its antecedents are
// printed
var TraceNotifications bool
var notificationGraphs []*NotificationGraph

// EnableNotificationPrint turns on the printing of notifications during tracing.
var EnableNotificationPrint bool

// traceFunctionCalls determines whether individual function calls will be graced. Its primary purpose is
// to understand what notifications resulted in the call to the function. When true, every time a function call
// is executed a graph of the function call and its antecedent notifications will be created and added to the
// functionCallGraphs
var traceFunctionCalls bool
var functionCallGraphs []*FunctionCallGraph

// notificationsLimit places an absolute limit on the number of notifications allowed. A value of 0 means no limit.
var notificationsLimit int

// notificationsCount counts the number of notifications that have occurred
var notificationsCount int

// ClearFunctionCallGraphs deletes all existing FunctionCallGraphs
func ClearFunctionCallGraphs() {
	functionCallGraphs = nil
}

// ClearNotificationGraphs deletes all existing NotificationGraphs
func ClearNotificationGraphs() {
	notificationGraphs = nil
}

// GetFunctionCallGraphs returns the array of FunctionCallGraphs that have been created
func GetFunctionCallGraphs() []*FunctionCallGraph {
	return functionCallGraphs
}

// GetNotificationGraphs returns the array of NotificationGraphs that have been created
func GetNotificationGraphs() []*NotificationGraph {
	return notificationGraphs
}

// GetNotificationsLimit returns the current limit on the number of notifications allowed. The default value
// of 0 indicates that there is no limit. This is the normal operating value. A limit may be imposed to aid
// in troubleshooting scenarios that have large numbers of notifications within a single "transaction" as
// defined as the scope of the changes that occur before locks are released
func GetNotificationsLimit() int {
	return notificationsLimit
}

// SetNotificationsLimit is provided as a debugging aid. It limits the number of change notifications allowed.
// A value of 0 is unlimited and is the normal production value.
func SetNotificationsLimit(limit int) {
	notificationsLimit = limit
	notificationsCount = 0
	notificationGraphs = nil
}
