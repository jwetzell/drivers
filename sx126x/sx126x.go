package sx126x

import (
	"time"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/internal/pin"
)

type Device struct {
	spi      drivers.SPI
	nssPin   pin.Output
	resetPin pin.Output
	busyPin  pin.Input
	spiTxBuf []byte
	spiRxBuf []byte
}

func New(spi drivers.SPI, nssPin pin.Output, resetPin pin.Output, busyPin pin.Input) *Device {
	return &Device{
		spi:      spi,
		nssPin:   nssPin,
		resetPin: resetPin,
		busyPin:  busyPin,
		spiTxBuf: make([]byte, 255), // TODO: optimize buffer size
		spiRxBuf: make([]byte, 255),
	}
}

func (d *Device) Reset() {
	// TODO(jwetzell): check timing requirements for this
	d.resetPin.Set(false)
	time.Sleep(10 * time.Millisecond)
	d.resetPin.Set(true)
	time.Sleep(10 * time.Millisecond)
}

func (d *Device) WaitWhileBusy() {
	// TODO(jwetzell): better way to do this?
	for {
		if !d.busyPin.Get() {
			return
		}
	}
}

func (d *Device) SetSleep(sleepConfig uint8) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_SLEEP, sleepConfig)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetStandby(standbyConfig uint8) {
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_STANDBY, standbyConfig)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetFs() {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_FS)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetTx(timeoutMs uint32) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	rtcSteps := timeoutMsToRtcSteps(timeoutMs)
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_TX, uint8((rtcSteps>>16)&0xFF), uint8((rtcSteps>>8)&0xFF), uint8(rtcSteps&0xFF))
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetRx(timeoutMs uint32) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	rtcSteps := timeoutMsToRtcSteps(timeoutMs)
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_RX, uint8((rtcSteps>>16)&0xFF), uint8((rtcSteps>>8)&0xFF), uint8(rtcSteps&0xFF))
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) StopTimerOnPreamble(enable bool) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_STOP_TIMER_ON_PREAMBLE)
	if enable {
		d.spiTxBuf = append(d.spiTxBuf, 1)
	} else {
		d.spiTxBuf = append(d.spiTxBuf, 0)
	}
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetRxDutyCycle(rxPeriod uint32, sleepPeriod uint32) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_RX_DUTY_CYCLE, uint8((rxPeriod>>16)&0xFF), uint8((rxPeriod>>8)&0xFF), uint8(rxPeriod&0xFF))
	d.spiTxBuf = append(d.spiTxBuf, uint8((sleepPeriod>>16)&0xFF), uint8((sleepPeriod>>8)&0xFF), uint8(sleepPeriod&0xFF))
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetCAD() {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_CAD)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetTxContinuousWave() {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_TX_CONTINUOUS_WAVE)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetTxInfinitePreamble() {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_TX_INFINITE_PREAMBLE)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetRegulatorMode(mode uint8) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_REGULATOR_MODE, mode)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) Calibrate(calibParam uint8) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_CALIBRATE, calibParam)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) CalibrateImage(freq1 uint8, freq2 uint8) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_CALIBRATE_IMAGE, freq1, freq2)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetPaConfig(paDutyCyle, hpMax, deviceSel uint8) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_PA_CONFIG, paDutyCyle, hpMax, deviceSel, PA_LUT)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetRxTxFallbackMode(fallbackMode uint8) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_RX_TX_FALLBACK_MODE, fallbackMode)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) WriteRegister(addr uint16, data []byte) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_WRITE_REGISTER, uint8((addr>>8)&0xFF), uint8(addr&0xFF))
	d.spiTxBuf = append(d.spiTxBuf, data...)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) ReadRegister(addr uint16) uint8 {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_READ_REGISTER, uint8((addr>>8)&0xFF), uint8(addr&0xFF), 0x00, 0x00)
	d.spiRxBuf = d.spiRxBuf[:5]
	d.spi.Tx(d.spiTxBuf, d.spiRxBuf)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
	return d.spiRxBuf[4]
}

func (d *Device) WriteBuffer(offset uint8, data []byte) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_WRITE_BUFFER, offset)
	d.spiTxBuf = append(d.spiTxBuf, data...)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) ReadBuffer(offset uint8, length uint8) []byte {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_READ_BUFFER, offset, 0x00)
	for i := uint8(0); i < length; i++ {
		d.spiTxBuf = append(d.spiTxBuf, 0x00)
	}
	d.spiRxBuf = d.spiRxBuf[:len(d.spiTxBuf)]
	d.spi.Tx(d.spiTxBuf, d.spiRxBuf)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
	return d.spiRxBuf[3:]
}

func (d *Device) SetDioIrqParams(irqMask, dio1Mask, dio2Mask, dio3Mask uint16) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_DIO_IRQ_PARAMS)
	d.spiTxBuf = append(d.spiTxBuf, uint8((irqMask>>8)&0xFF), uint8(irqMask&0xFF))
	d.spiTxBuf = append(d.spiTxBuf, uint8((dio1Mask>>8)&0xFF), uint8(dio1Mask&0xFF))
	d.spiTxBuf = append(d.spiTxBuf, uint8((dio2Mask>>8)&0xFF), uint8(dio2Mask&0xFF))
	d.spiTxBuf = append(d.spiTxBuf, uint8((dio3Mask>>8)&0xFF), uint8(dio3Mask&0xFF))
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) GetIrqStatus() uint16 {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_GET_IRQ_STATUS, 0x00, 0x00, 0x00)
	d.spiRxBuf = d.spiRxBuf[:4]
	d.spi.Tx(d.spiTxBuf, d.spiRxBuf)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
	return uint16(d.spiRxBuf[2])<<8 | uint16(d.spiRxBuf[3])
}

func (d *Device) ClearIrqStatus(irqMask uint16) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_CLEAR_IRQ_STATUS, uint8((irqMask>>8)&0xFF), uint8(irqMask&0xFF))
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetDIO2AsRfSwitchCtrl(enable bool) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_DIO2_AS_RF_SWITCH_CTRL)
	if enable {
		d.spiTxBuf = append(d.spiTxBuf, 1)
	} else {
		d.spiTxBuf = append(d.spiTxBuf, 0)
	}
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetDIO3AsTCXOCtrl(tcxoVoltage uint8, delayUs uint32) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	delay := uint32(float32(delayUs) / 15.625)
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_DIO3_AS_TCXO_CTRL, tcxoVoltage)
	d.spiTxBuf = append(d.spiTxBuf, uint8((delay>>16)&0xFF), uint8((delay>>8)&0xFF), uint8(delay&0xFF))
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetRfFrequency(frequency uint32) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	convertedFrequency := uint32((uint64(frequency) << 25) / 32000000)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_RF_FREQUENCY, uint8((convertedFrequency>>24)&0xFF), uint8((convertedFrequency>>16)&0xFF), uint8((convertedFrequency>>8)&0xFF), uint8(convertedFrequency&0xFF))
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetPacketType(packetType uint8) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_PACKET_TYPE, packetType)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) GetPacketType() uint8 {
	// TODO(jwetzell): check the NOP part of this
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_GET_PACKET_TYPE, 0x00, 0x00)
	d.spiRxBuf = d.spiRxBuf[:3]
	d.spi.Tx(d.spiTxBuf, d.spiRxBuf)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
	return d.spiRxBuf[2]
}

func (d *Device) SetTxParams(power int8, rampTime uint8) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_TX_PARAMS, uint8(power), rampTime)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

// GFSK: param1-3 = bitrate, param4 = pulseShape, param5 = bandwidth, param6-8 = fDev
func (d *Device) SetModulationParamsGFSK(bitrate uint32, pulseShape uint8, bandwidth uint8, fDev uint32) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_MODULATION_PARAMS)
	d.spiTxBuf = append(d.spiTxBuf, uint8((bitrate>>16)&0xFF), uint8((bitrate>>8)&0xFF), uint8(bitrate&0xFF))
	d.spiTxBuf = append(d.spiTxBuf, pulseShape, bandwidth)
	d.spiTxBuf = append(d.spiTxBuf, uint8((fDev>>16)&0xFF), uint8((fDev>>8)&0xFF), uint8(fDev&0xFF))
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetModulationParamsLoRa(spreadingFactor, bandwidth, codingRate, lowDataRateOptimize uint8) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_MODULATION_PARAMS, spreadingFactor, bandwidth, codingRate, lowDataRateOptimize)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetPacketParamsGFSK(preambleLength uint16, preambleDetectorLength, syncWordLength, addressFiltering, packetType, payloadLength, crcType, whitening uint8) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_PACKET_PARAMS)
	d.spiTxBuf = append(d.spiTxBuf, uint8((preambleLength>>8)&0xFF), uint8(preambleLength&0xFF))
	d.spiTxBuf = append(d.spiTxBuf, preambleDetectorLength, syncWordLength, addressFiltering, packetType, payloadLength, crcType, whitening)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetPacketParamsLoRa(preambleLength uint16, headerType, payloadLength, crcType, invertIQ uint8) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_PACKET_PARAMS)
	d.spiTxBuf = append(d.spiTxBuf, uint8((preambleLength>>8)&0xFF), uint8(preambleLength&0xFF))
	d.spiTxBuf = append(d.spiTxBuf, headerType, payloadLength, crcType, invertIQ)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetCADParams(cadSymbolNum uint8, cadDetPeak uint8, cadDetMin uint8, cadExitMode uint8, cadTimeoutMs uint32) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_CAD_PARAMS)
	d.spiTxBuf = append(d.spiTxBuf, cadSymbolNum, cadDetPeak, cadDetMin, cadExitMode)
	rtcSteps := timeoutMsToRtcSteps(cadTimeoutMs)
	d.spiTxBuf = append(d.spiTxBuf, uint8((rtcSteps>>16)&0xFF), uint8((rtcSteps>>8)&0xFF), uint8(rtcSteps&0xFF))
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetBufferBaseAddress(txBase uint8, rxBase uint8) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_BUFFER_BASE_ADDRESS, txBase, rxBase)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) SetLoraSymbNumTimeout(symbNum uint8) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_SET_LORA_SYMB_NUM_TIMEOUT, symbNum)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) GetStatus() (uint8, uint8) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_GET_STATUS, 0x00)
	d.spiRxBuf = d.spiRxBuf[:2]
	d.spi.Tx(d.spiTxBuf, d.spiRxBuf)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
	status := d.spiRxBuf[1]
	chipMode := (status & CHIP_MODE_MASK) >> 4
	commandStatus := (status & COMMAND_STATUS_MASK) >> 1
	return chipMode, commandStatus
}

func (d *Device) GetRxBufferStatus() (payloadLength uint8, startBufferPointer uint8) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_GET_RX_BUFFER_STATUS, 0x00, 0x00, 0x00)
	d.spiRxBuf = d.spiRxBuf[:4]
	d.spi.Tx(d.spiTxBuf, d.spiRxBuf)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
	return d.spiRxBuf[2], d.spiRxBuf[3]
}

func (d *Device) GetPacketStatus() (uint8, uint8, uint8) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_GET_PACKET_STATUS, 0x00, 0x00, 0x00, 0x00)
	d.spiRxBuf = d.spiRxBuf[:5]
	d.spi.Tx(d.spiTxBuf, d.spiRxBuf)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
	return d.spiRxBuf[2], d.spiRxBuf[3], d.spiRxBuf[4]
}

func (d *Device) GetRssiInst() uint8 {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_GET_RSSI_INST, 0x00, 0x00)
	d.spiRxBuf = d.spiRxBuf[:3]
	d.spi.Tx(d.spiTxBuf, d.spiRxBuf)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
	return d.spiRxBuf[2]
}

func (d *Device) GetStats() (uint16, uint16, uint16) {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_GET_STATS, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00)
	d.spiRxBuf = d.spiRxBuf[:8]
	d.spi.Tx(d.spiTxBuf, d.spiRxBuf)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
	return uint16(d.spiRxBuf[2])<<8 | uint16(d.spiRxBuf[3]), uint16(d.spiRxBuf[4])<<8 | uint16(d.spiRxBuf[5]), uint16(d.spiRxBuf[6])<<8 | uint16(d.spiRxBuf[7])
}

func (d *Device) ResetStats() {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_RESET_STATS, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func (d *Device) GetDeviceErrors() uint16 {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_GET_DEVICE_ERRORS, 0x00, 0x00, 0x00)
	d.spiRxBuf = d.spiRxBuf[:4]
	d.spi.Tx(d.spiTxBuf, d.spiRxBuf)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
	return uint16(d.spiRxBuf[2])<<8 | uint16(d.spiRxBuf[3])
}

func (d *Device) ClearDeviceErrors() {
	d.WaitWhileBusy()
	d.nssPin.Set(false)
	d.spiTxBuf = d.spiTxBuf[:0]
	d.spiTxBuf = append(d.spiTxBuf, CMD_CLEAR_DEVICE_ERRORS, 0x00, 0x00)
	d.spi.Tx(d.spiTxBuf, nil)
	d.nssPin.Set(true)
	d.WaitWhileBusy()
}

func timeoutMsToRtcSteps(timeoutMs uint32) uint32 {
	r := uint32(float32(timeoutMs*1000) / 15.625)
	return r
}
