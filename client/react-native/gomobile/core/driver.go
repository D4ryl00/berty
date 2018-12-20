package core

type NativeLogger interface {
	Log(level, namespace, message string) error
	LevelEnabler(level string) bool
}

type NativeNotification interface {
	DisplayNativeNotification(title, body, icon, sound string) error
}
