package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/cs104"
)

const interrogationResolution = 10 * time.Minute // 总召唤每间隔15分钟进行一次
//const counterInterrogationResolution = 15 * time.Minute // 电能召唤每间隔15分钟进行一次
const clockSyncResolution = 20 * time.Minute // 总召唤每间隔15分钟进行一次

//var interrogationTicker = time.NewTicker(interrogationResolution) // 总召唤

const fileRoot = "./raw_folder/"
const FILE_PACKET_SIZE = 220

var fileMapping = make(map[uint32]uint32, 100)
var fileBuffer = make(map[uint32][]byte, 100)
var fileNames = make(map[uint32]string, 100)
var fileSegment = make(map[uint32]uint32, 100)

func main() {
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
			err := asdu.InterrogationCmd(c, asdu.CauseOfTransmission{Cause: asdu.Activation}, asdu.CommonAddr(1), asdu.QOIStation)
			if err != nil {
				log.Println("InterrogationCmd falied", err, c, asduPack)
				return
			} else {
				log.Println("InterrogationCmd success", err, c, asduPack)
			}
			//err = asdu.CounterInterrogationCmd(c, asdu.CauseOfTransmission{Cause: asdu.Activation}, asdu.CommonAddr(1), asdu.QualifierCountCall{asdu.QCCTotal, asdu.QCCFrzRead})
			//if err != nil {
			//	log.Println("CounterInterrogationCmd falied", err, c, asduPack)
			//} else {
			//	log.Println("CounterInterrogationCmd success", err, c, asduPack)
			//}
			time.Sleep(interrogationResolution)
		}
	}()
	go func() {
		for {
			err := asdu.ClockSynchronizationCmd(c, asdu.CauseOfTransmission{Cause: asdu.Activation}, asdu.CommonAddr(1), time.Now())
			if err != nil {
				log.Println("ClockSynchronizationCmd falied", err, c, asduPack)
				return
			} else {
				log.Println("ClockSynchronizationCmd success", err, c, asduPack)
			}
			time.Sleep(clockSyncResolution)
		}
	}()
	return nil
	//return asdu.InterrogationCmd(c, asdu.CauseOfTransmission{Cause: asdu.Activation}, asdu.CommonAddr(1), asdu.QOIStation)
}

func (sf *mysrv) InterrogationHandler(c asdu.Connect, asduPack *asdu.ASDU, qoi asdu.QualifierOfInterrogation) error {
	log.Println("InterrogationHandler", qoi)
	//asduPack.SendReplyMirror(c, asdu.ActivationTerm)
	return nil
	//return asdu.CounterInterrogationCmd(c, asdu.CauseOfTransmission{Cause: asdu.Activation}, asdu.CommonAddr(1), asdu.QualifierCountCall{asdu.QCCTotal, asdu.QCCFrzRead})
}
func (sf *mysrv) CounterInterrogationHandler(c asdu.Connect, asduPack *asdu.ASDU, qcc asdu.QualifierCountCall) error {
	log.Println("CounterInterrogationHandler", asduPack)
	return nil
	//return asdu.ClockSynchronizationCmd(c, asdu.CauseOfTransmission{Cause: asdu.Activation}, asdu.CommonAddr(1), time.Now())
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
		log.Println("遥信单点 GetSinglePoint", len(ret))
		for i, v := range ret {
			log.Println(c.ServerId(), i, v.Ioa, v.Value, v.Time)
		}
		log.Println("========")
		break
	case asdu.M_DP_TB_1: // 3.2 遥信及SOE事件上报上送 带 CP56Time2a 时标的双点信息
	case asdu.M_DP_NA_1: // 3.2 遥信及SOE事件上报上送 双点信息
		var ret = asduPack.GetDoublePoint()
		log.Println("遥信双点 GetDoublePoint", len(ret))
		for i, v := range ret {
			log.Println(c.ServerId(), i, v.Ioa, v.Value, v.Time)
		}
		log.Println("========")
		break
	case asdu.M_ME_NC_1: // 3.1 遥测上送 短浮点数
		var ret = asduPack.GetMeasuredValueFloat()
		log.Println("遥测上送 GetMeasuredValueFloat", len(ret))
		for i, v := range ret {
			log.Println(c.ServerId(), i, v.Ioa, v.Value, v.Time)
		}
		log.Println("========")
		break
	case asdu.M_IT_NB_1: // 3.3.2 电能量上送 累计量，浮点短数
		var ret = asduPack.GetIntegratedFloatTotals()
		log.Println("电能量上送 累计量 GetIntegratedFloatTotals", len(ret))
		for i, v := range ret {
			log.Println(c.ServerId(), i, v.Ioa, v.Value, v.Time)
		}
		log.Println("========")
		break
	case asdu.M_ME_NA_1: // 测量值，归一化值
		var ret = asduPack.GetMeasuredValueNormal()
		log.Println("测量值，归一化值 GetMeasuredValueNormal", len(ret))
		for i, v := range ret {
			log.Println(c.ServerId(), i, v.Ioa, v.Value, v.Time)
		}
		log.Println("========")
		break
	case asdu.F_NF_NA_1: // 新文件主动上报
		asduPack.DecodeInfoObjAddr()
		extraPacketType := asduPack.DecodeByte()
		opFlag := asduPack.DecodeByte()
		fileNameSize := int(asduPack.DecodeByte())
		fileName := asduPack.DecodeString(fileNameSize)
		fileId := asduPack.DecodeBitsString32()
		//dSn := asduPack.DecodeUint16()
		//dActionCount := asduPack.DecodeUint16()
		log.Printf("新文件主动上报 附加数据包类型 %d 操作标识 %d 文件名长度 %d 文件名 %s fileId %d\n", extraPacketType, opFlag, fileNameSize,
			fileName, fileId)
		// get new file upload
		// send 读文件召唤报文
		e := asdu.FileReadActiveCmd(c, asdu.CauseOfTransmission{Cause: asdu.Activation}, asdu.CommonAddr(1), fileName)
		log.Println("新文件主动上报 读文件召唤报文发送", e)
		break
	case asdu.F_FR_NA_1:
		asduPack.DecodeInfoObjAddr()
		extraPacketType := asduPack.DecodeByte()                                                   // 附加数据包类型
		opFlag := asduPack.DecodeByte()                                                            // 操作标识
		if asduPack.Coa.Cause == asdu.ActivationCon && extraPacketType == 0x02 && opFlag == 0x04 { // 激活确认报文
			actRet := asduPack.DecodeByte() // 结果描述字
			fileNameSize := int(asduPack.DecodeByte())
			fileName := asduPack.DecodeString(fileNameSize)
			fileId := asduPack.DecodeBitsString32()
			fileContentSize := asduPack.DecodeBitsString32()
			_, has := fileMapping[fileId]
			if actRet == 0 && !has {
				// 成功
				log.Println("收到读文件激活确认成功报文", fileContentSize, fileName, fileId)
				fileMapping[fileId] = fileContentSize
				fileNames[fileId] = fileName
				fileBuffer[fileId] = make([]byte, fileContentSize)
				fileSegment[fileId] = 0
			} else {
				// 失败
				log.Println("收到读文件激活确认失败报文", fileContentSize, fileName, fileId)
			}
		} else if asduPack.Coa.Cause == asdu.Request && extraPacketType == 0x02 && opFlag == 0x05 { // 读文件数据报文
			fileId := asduPack.DecodeBitsString32()
			segmentId := asduPack.DecodeBitsString32() // 数据段号,可以使用文件内容的偏移指针值
			cSegmentId, hasSegmentId := fileSegment[fileId]
			notEOF := asduPack.DecodeByte()
			fSize, has := fileMapping[fileId]
			if has && hasSegmentId && cSegmentId == segmentId {
				cSize := fSize - segmentId
				if cSize > FILE_PACKET_SIZE {
					cSize = FILE_PACKET_SIZE
				}
				contents := asduPack.DecodeBytes(int(cSize))
				chk := asduPack.DecodeByte()
				if contents != nil && Checksum(contents) == chk && len(fileBuffer[fileId]) > int(segmentId) {
					fmt.Println("segmentId", segmentId, len(fileBuffer[fileId]))
					copy(fileBuffer[fileId][segmentId:], contents)
					log.Println("收到文件传输帧成功报文", fileId, segmentId, cSize, notEOF, "RECV", len(fileBuffer[fileId]))
					e := asdu.FileReadRequestCmd(c, asdu.CauseOfTransmission{Cause: asdu.Request}, asdu.CommonAddr(1), fileId, segmentId, false)
					fileSegment[fileId] += cSize
					log.Println("文件传输帧成功回复报文发送", e)
					if notEOF == 0 {
						log.Println("文件接受完毕，写出", fileRoot, fileNames[fileId])
						prefix := strings.Split(fileNames[fileId], "_")[0]
						_ = os.MkdirAll(filepath.Join(fileRoot, prefix), os.ModePerm)
						err := ioutil.WriteFile(filepath.Join(fileRoot, prefix, fileNames[fileId]), fileBuffer[fileId], 0644)
						if err != nil {
							log.Println("文件接受完毕，写出失败", fileRoot+fileNames[fileId], err)
						}
						delete(fileMapping, fileId)
						delete(fileNames, fileId)
						delete(fileBuffer, fileId)
						delete(fileSegment, fileId)
					}
					break
				}
				log.Println("收到文件传输帧成功报文", fileId, segmentId, cSize, notEOF, "RECV", len(fileBuffer[fileId]), Checksum(contents) == chk, len(fileBuffer[fileId]) > int(segmentId), len(fileBuffer[fileId]), int(segmentId))
			}
			log.Println("收到文件传输帧成功报文", fileId, fSize, cSegmentId, segmentId, notEOF, "RECV", len(fileBuffer[fileId]), len(fileBuffer[fileId]) > int(segmentId), len(fileBuffer[fileId]), int(segmentId))
			e := asdu.FileReadRequestCmd(c, asdu.CauseOfTransmission{Cause: asdu.Request}, asdu.CommonAddr(1), fileId, segmentId, true)
			log.Println("文件传输帧成功回复报文发送[failed]", e)
		}
		break
	default:
		log.Println("unknown type", asduPack.String())
	}
	return nil
}

func Checksum(data []byte) byte {
	sum := 0
	for _, value := range data {
		sum += int(value)
	}
	return byte(sum)
}
