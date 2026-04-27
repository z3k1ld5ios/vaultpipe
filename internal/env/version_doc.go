// Package env — Versioner
//
// Versioner tracks successive states of an environment map using a
// monotonically increasing generation counter combined with a short
// content checksum.
//
// The generation counter increments only when the content of the map
// actually changes, making it safe to use as a cache-invalidation key
// or a lightweight change-detection signal without storing the full
// previous state at the call site.
//
// Usage:
//
//	v := env.NewVersioner(initialSecrets)
//	ver := v.Apply(updatedSecrets)
//	fmt.Println(ver) // gen=1 chk=a3f9...
package env
