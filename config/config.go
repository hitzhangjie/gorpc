package config

type Config interface {
	Read(key string, dftValue string) string
	ReadInt(key string, dftValue int) int
	ReadBool(key string, dftValue bool) bool
}
