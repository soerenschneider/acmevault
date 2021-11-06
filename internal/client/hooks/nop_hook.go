package hooks

// NopHook does nothing and is used if no (valid) hook provided was given.
type NopHook struct {
}

func (hook *NopHook) Invoke() error {
	return nil
}
