package astits

import (
	"fmt"
	"time"

	"github.com/asticode/go-astikit"
)

// parseDVBTime parses a DVB time
// This field is coded as 16 bits giving the 16 LSBs of MJD followed by 24 bits coded as 6 digits in 4 - bit Binary
// Coded Decimal (BCD). If the start time is undefined (e.g. for an event in a NVOD reference service) all bits of the
// field are set to "1".
// I apologize for the computation which is really messy but details are given in the documentation
// Page: 160 | Annex C | Link: https://www.dvb.org/resources/public/standards/a38_dvb-si_specification.pdf
// (barbashov) the link above can be broken, alternative: https://dvb.org/wp-content/uploads/2019/12/a038_tm1217r37_en300468v1_17_1_-_rev-134_-_si_specification.pdf
func parseDVBTime(i *astikit.BytesIterator) (t time.Time, err error) {
	// Get next 2 bytes
	var bs []byte
	if bs, err = i.NextBytesNoCopy(2); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}

	// Date
	var mjd = uint16(bs[0])<<8 | uint16(bs[1])
	var yt = int((float64(mjd) - 15078.2) / 365.25)
	var mt = int((float64(mjd) - 14956.1 - float64(int(float64(yt)*365.25))) / 30.6001)
	var d = int(float64(mjd) - 14956 - float64(int(float64(yt)*365.25)) - float64(int(float64(mt)*30.6001)))
	var k int
	if mt == 14 || mt == 15 {
		k = 1
	}
	var y = yt + k
	var m = mt - 1 - k*12
	t, _ = time.Parse("06-01-02", fmt.Sprintf("%d-%d-%d", y, m, d))

	// Time
	var s time.Duration
	if s, err = parseDVBDurationSeconds(i); err != nil {
		err = fmt.Errorf("astits: parsing DVB duration seconds failed: %w", err)
		return
	}
	t = t.Add(s)
	return
}

// parseDVBDurationMinutes parses a minutes duration
// 16 bit field containing the duration of the event in hours, minutes. format: 4 digits, 4 - bit BCD = 18 bit
func parseDVBDurationMinutes(i *astikit.BytesIterator) (d time.Duration, err error) {
	var bs []byte
	if bs, err = i.NextBytesNoCopy(2); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}
	d = parseDVBDurationByte(bs[0])*time.Hour + parseDVBDurationByte(bs[1])*time.Minute
	return
}

// parseDVBDurationSeconds parses a seconds duration
// 24 bit field containing the duration of the event in hours, minutes, seconds. format: 6 digits, 4 - bit BCD = 24 bit
func parseDVBDurationSeconds(i *astikit.BytesIterator) (d time.Duration, err error) {
	var bs []byte
	if bs, err = i.NextBytesNoCopy(3); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}
	d = parseDVBDurationByte(bs[0])*time.Hour + parseDVBDurationByte(bs[1])*time.Minute + parseDVBDurationByte(bs[2])*time.Second
	return
}

// parseDVBDurationByte parses a duration byte
func parseDVBDurationByte(i byte) time.Duration {
	return time.Duration(uint8(i)>>4*10 + uint8(i)&0xf)
}

func writeDVBTime(w *astikit.BitsWriter, t time.Time) (int, error) {
	year := t.Year() - 1900
	month := t.Month()
	day := t.Day()

	l := 0
	if month <= time.February {
		l = 1
	}

	mjd := 14956 + day + int(float64(year-l)*365.25) + int(float64(int(month)+1+l*12)*30.6001)

	d := t.Sub(t.Truncate(24 * time.Hour))

	b := astikit.NewBitsWriterBatch(w)

	b.Write(uint16(mjd))
	bytesWritten, err := writeDVBDurationSeconds(w, d)
	if err != nil {
		return 2, err
	}

	return bytesWritten + 2, b.Err()
}

func writeDVBDurationMinutes(w *astikit.BitsWriter, d time.Duration) (int, error) {
	b := astikit.NewBitsWriterBatch(w)

	hours := uint8(d.Hours())
	minutes := uint8(int(d.Minutes()) % 60)

	b.Write(dvbDurationByteRepresentation(hours))
	b.Write(dvbDurationByteRepresentation(minutes))

	return 2, b.Err()
}

func writeDVBDurationSeconds(w *astikit.BitsWriter, d time.Duration) (int, error) {
	b := astikit.NewBitsWriterBatch(w)

	hours := uint8(d.Hours())
	minutes := uint8(int(d.Minutes()) % 60)
	seconds := uint8(int(d.Seconds()) % 60)

	b.Write(dvbDurationByteRepresentation(hours))
	b.Write(dvbDurationByteRepresentation(minutes))
	b.Write(dvbDurationByteRepresentation(seconds))

	return 3, b.Err()
}

func dvbDurationByteRepresentation(n uint8) uint8 {
	return (n/10)<<4 | n%10
}
