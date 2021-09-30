package distributor

type DriverInterface interface {
	//Migration миграция изменений.
	Migration() (err error)
	//CanHandleEvent может ли бот обработать событие?
	CanHandleEvent(eventID, message string) (ok bool, msg string, errs error)
}
