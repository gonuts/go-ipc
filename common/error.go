package common

type PluginError struct {
	msg string
}

func (e *PluginError) Error() string {
	return e.msg
}

func NewPluginError(message string) *PluginError {
	err := &PluginError{message}
	return err
}
