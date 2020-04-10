package tableflux

type tableFluxImpl = func(input string) (bool, string, string)

var impl tableFluxImpl

// For use by the CGo module. It registers itself when compiled in. Doesn't
// need to be public, except to satisfy staticcheck.
func RegisterTableFlux(tableflux tableFluxImpl) {
	impl = tableflux
}

// Is it even compiled in?
func Enabled() bool {
	return impl != nil
}

// Translate tableflux to flux. Returns (ok, flux, log)
//   ok:   Comes back true if the transformation succeeded, false otherwise
//   flux: The result of the transformation. This may contain a partial result
//         when the transformation fails.
//   log:  Log information. This may contain log information regardless of
//         success. If the transormation fails it will also contain the reason.
func TableFlux(input string) (bool, string, string) {
	if impl == nil {
		return false, "", "error: TableFlux not compiled in"
	}

	return impl(input)
}
