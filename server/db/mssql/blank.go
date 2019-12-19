// +build !mssql

// This file is needed for conditional compilation. It's used when
// the build tag 'mssql' is not defined. Otherwise the adapter.go
// is compiled.

package mssql
