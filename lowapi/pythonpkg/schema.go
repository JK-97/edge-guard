package pythonpkg

type PkgClient struct {
	InternalPkg []PythonPkg `json:"internal_pkg"`
	AllPkg      []PythonPkg `json:"all_pkg"`
}

type PythonPkg struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
