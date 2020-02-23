package config

// Config config
//
// Read, read string value by `key`, if not found, return dftValue
// ReadInt, read int value by `key`, if not found, return dftValue
// ReadBool, read bool value by `key`, if not found, return dftValue
type Config interface {
	Read(key string, dftValue string) string
	ReadInt(key string, dftValue int) int
	ReadBool(key string, dftValue bool) bool
}
