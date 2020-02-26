package main

import (
	"log"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/cs104"
)

const interrogationResolution = 15 * time.Minute // 总召唤每间隔15分钟进行一次

var interrogationTicker = time.NewTicker(interrogationResolution) // 总召唤

func main() {
	defer func() {
		interrogationTicker.Stop()
	}()
	srv := cs104.NewServer(&mysrv{})
	srv.LogMode(true)
	srv.ListenAndServer(":2404")
}

type mysrv struct{}

func (sf *mysrv) EndOfInitializationHandler(c asdu.Connect, asduPack *asdu.ASDU, iaddr asdu.InfoObjAddr, coi asdu.CauseOfInitial) error {
	log.Println("EndOfInitializationHandler", iaddr, coi)
	//asduPack.SendReplyMirror(c, asdu.ActivationCon)
	go func() {
		for {
			time.Sleep(interrogationResolution)
			err := asdu.InterrogationCmd(c, asdu.CauseOfTransmission{Cause: asdu.Activation}, asdu.CommonAddr(1), asdu.QOIStation)
			if err != nil {
				log.Println("InterrogationCmd falied", err, c, asduPack)
				return
			} else {
				log.Println("InterrogationCmd success", err, c, asduPack)
			}
		}
	}()
	return asdu.InterrogationCmd(c, asdu.CauseOfTransmission{Cause: asdu.Activation}, asdu.CommonAddr(1), asdu.QOIStation)
}

func (sf *mysrv) InterrogationHandler(c asdu.Connect, asduPack *asdu.ASDU, qoi asdu.QualifierOfInterrogation) error {
	log.Println("InterrogationHandler", qoi)
	//asduPack.SendReplyMirror(c, asdu.ActivationTerm)
	return asdu.ClockSynchronizationCmd(c, asdu.CauseOfTransmission{Cause: asdu.Activation}, asdu.CommonAddr(1), time.Now())
}
func (sf *mysrv) CounterInterrogationHandler(c asdu.Connect, asduPack *asdu.ASDU, qcc asdu.QualifierCountCall) error {
	log.Println("CounterInterrogationHandler", asduPack)
	return nil
}
func (sf *mysrv) ReadHandler(asdu.Connect, *asdu.ASDU, asdu.InfoObjAddr) error {
	log.Println("ReadHandler")
	return nil
}
func (sf *mysrv) ClockSyncHandler(asdu.Connect, *asdu.ASDU, time.Time) error {
	log.Println("ClockSyncHandler")
	return nil
}
func (sf *mysrv) ResetProcessHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierOfResetProcessCmd) error {
	log.Println("ResetProcessHandler")
	return nil
}
func (sf *mysrv) DelayAcquisitionHandler(asdu.Connect, *asdu.ASDU, uint16) error {
	log.Println("DelayAcquisitionHandler")
	return nil
}
func (sf *mysrv) ASDUHandler(c asdu.Connect, asduPack *asdu.ASDU) error {
	log.Println("ASDUHandler", asduPack.String(), asduPack)
	switch asduPack.Type {
	case asdu.M_SP_TB_1: // 3.2 遥信及SOE事件上报上送 带 CP56Time2a 时标的单点信息
	case asdu.M_SP_NA_1: // 3.2 遥信及SOE事件上报上送 单点信息
		var ret = asduPack.GetSinglePoint()
		log.Println("GetSinglePoint", len(ret), ret)
		break
	case asdu.M_DP_TB_1: // 3.2 遥信及SOE事件上报上送 带 CP56Time2a 时标的双点信息
	case asdu.M_DP_NA_1: // 3.2 遥信及SOE事件上报上送 双点信息
		var ret = asduPack.GetDoublePoint()
		log.Println("GetDoublePoint", len(ret), ret)
		break
	case asdu.M_ME_NC_1: // 3.1 遥测上送 短浮点数
		var ret = asduPack.GetMeasuredValueFloat()
		log.Println("GetMeasuredValueFloat", len(ret), ret)
		break
	case asdu.M_IT_NB_1: // 3.3.2 电能量上送 累计量，浮点短数
		var ret = asduPack.GetIntegratedFloatTotals()
		log.Println("GetIntegratedFloatTotals", len(ret), ret)
		break
	case asdu.M_ME_NA_1: // 测量值，归一化值
		var ret = asduPack.GetMeasuredValueNormal()
		log.Println("GetMeasuredValueNormal", len(ret), ret)
		break
	default:
		log.Println("unknown type", asduPack.String())
	}
	return nil
}
