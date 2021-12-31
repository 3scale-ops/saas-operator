package resources

var (
	// DefaultExcludedPaths is a list of jsonpaths paths to ignore during reconciliation
	DefaultExcludedPaths []string = []string{
		"/metadata/creationTimestamp",
		"/metadata/deletionGracePeriodSeconds",
		"/metadata/deletionTimestamp",
		"/metadata/finalizers",
		"/metadata/generateName",
		"/metadata/generation",
		"/metadata/managedFields",
		"/metadata/ownerReferences",
		"/metadata/resourceVersion",
		"/metadata/selfLink",
		"/metadata/uid",
		"/status",
	}
)
