// Package filter provides allow/deny filtering of secret key-value maps
// before they are injected into a process environment.
//
// Use Filter with a Config to restrict which secrets are exposed,
// reducing the blast radius of unintended secret leakage.
package filter
