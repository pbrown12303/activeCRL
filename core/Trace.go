package core

// TraceChange is a boolean that, when true, turns on the tracing of ChangeNotifications
var TraceChange bool

// OmitHousekeepingCalls is a boolean that indicates whether housekeeping calls should be omitted when TraceChanve is enabled.
var OmitHousekeepingCalls bool

// OmitManageTreeNodesCalls is a boolean that indicates whether tree node management calls should be omitted when TraceChanve is enabled.
var OmitManageTreeNodesCalls bool

// OmitDiagramRelatedCalls is a boolean that indicates whether diagram management calls should be omitted when TraceChanve is enabled.
var OmitDiagramRelatedCalls bool

// TraceLocks is a boolean that, when true, turns on the tracing of locks
var TraceLocks = false

var notificationGraphs []*NotificationGraph

// EnableNotificationPrint turns on the printing of notifications during tracing.
var EnableNotificationPrint bool

// traceFunctionCalls determines whether individual function calls will be traced. Its primary purpose is
// to understand what notifications resulted in the call to the function. When true, every time a function call
// is executed a graph of the function call and its antecedent notifications will be created and added to the
// functionCallGraphs
// var traceFunctionCalls bool
var functionCallGraphs []*FunctionCallGraph

// notificationsLimit places an absolute limit on the number of notifications allowed. A value of 0 means no limit.
var notificationsLimit int

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
