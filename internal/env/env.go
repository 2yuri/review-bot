package env

import (
	"errors"
	"os"
	"strconv"
)

var errNilEnv = errors.New("env is empty")

type Env string

func Get(key string) Env {
	return Env(os.Getenv(key))
}

func (e Env) empty() bool {
	return e == ""
}

func (e Env) Int() (int, error) {
	if e.empty() {
		return 0, errNilEnv
	}

	return strconv.Atoi(string(e))
}

func (e Env) IntFallback(fallback int) int {
	if e.empty() {
		return fallback
	}

	v, err := strconv.Atoi(string(e))
	if err != nil {
		return fallback
	}

	return v
}

func (e Env) Int64() (int64, error) {
	if e.empty() {
		return 0, errNilEnv
	}

	return strconv.ParseInt(string(e), 10, 64)
}

func (e Env) Uint64Fallback(fallback uint64) uint64 {
	if e.empty() {
		return fallback
	}

	v, err := strconv.ParseUint(string(e), 10, 64)
	if err != nil {
		return fallback
	}

	return v
}

func (e Env) Int64Fallback(fallback int64) int64 {
	if e.empty() {
		return fallback
	}

	v, err := strconv.ParseInt(string(e), 10, 64)
	if err != nil {
		return fallback
	}

	return v
}

func (e Env) Float() (float64, error) {
	if e.empty() {
		return 0, errNilEnv
	}

	return strconv.ParseFloat(string(e), 64)
}

func (e Env) FloatFallback(fallback float64) float64 {
	if e.empty() {
		return fallback
	}

	v, err := strconv.ParseFloat(string(e), 64)
	if err != nil {
		return fallback
	}

	return v
}

func (e Env) Bool() (bool, error) {
	if e.empty() {
		return false, errNilEnv
	}

	return strconv.ParseBool(string(e))
}

func (e Env) BoolFallback(fallback bool) bool {
	if e.empty() {
		return fallback
	}

	v, err := strconv.ParseBool(string(e))
	if err != nil {
		return fallback
	}

	return v
}

func (e Env) String() string {
	return string(e)
}

func (e Env) StringFallback(fallback string) string {
	if e.empty() {
		return fallback
	}

	return string(e)
}
