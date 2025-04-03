//go:build darwin

package storage

// KiloByte determines how many bytes a KiloByte is.
// Since macOS 10.6 (Snow Leopard) Apple has used decimal (base-10) definition for user-facing interfaces.
const KiloByte = 1000
