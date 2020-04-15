package cs104

import (
	"time"

	"github.com/tsenc/go-iec104/asdu"
)

// ServerHandlerInterface is the interface of server handler
type ServerHandlerInterface interface {
	EndOfInitializationHandler(asdu.Connect, *asdu.ASDU, asdu.InfoObjAddr, asdu.CauseOfInitial) error
	InterrogationHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierOfInterrogation) error
	CounterInterrogationHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierCountCall) error
	ReadHandler(asdu.Connect, *asdu.ASDU, asdu.InfoObjAddr) error
	ClockSyncHandler(asdu.Connect, *asdu.ASDU, time.Time) error
	ResetProcessHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierOfResetProcessCmd) error
	DelayAcquisitionHandler(asdu.Connect, *asdu.ASDU, uint16) error
	ASDUHandler(asdu.Connect, *asdu.ASDU) error
}

// ClientHandlerInterface  is the interface of client handler
type ClientHandlerInterface interface {
	InterrogationHandler(asdu.Connect, *asdu.ASDU) error
	CounterInterrogationHandler(asdu.Connect, *asdu.ASDU) error
	ReadHandler(asdu.Connect, *asdu.ASDU) error
	TestCommandHandler(asdu.Connect, *asdu.ASDU) error
	ClockSyncHandler(asdu.Connect, *asdu.ASDU) error
	ResetProcessHandler(asdu.Connect, *asdu.ASDU) error
	DelayAcquisitionHandler(asdu.Connect, *asdu.ASDU) error
	ASDUHandler(asdu.Connect, *asdu.ASDU) error
}
