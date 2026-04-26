// Package env provides utilities for constructing, transforming, and auditing
// process environment variables sourced from HashiCorp Vault secrets.
//
// The audit_bridge module bridges raw env map comparisons into structured
// AuditSummary values, enabling downstream audit loggers and rotation hooks
// to record precise change sets without exposing secret values.
package env
