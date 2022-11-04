package transforms

import "math"

// DCT1DFast64 function returns result of DCT-II.
// DCT type II, unscaled. Algorithm by Byeong Gi Lee, 1984.
// Static implementation by Evan Oberholster, 2022.
func DCT1DFast64(input []float64) {
	var temp [64]float64
	for i := 0; i < 32; i++ {
		x, y := input[i], input[63-i]
		temp[i] = x + y
		temp[i+32] = (x - y) / dct64[i]
	}
	forwardTransformStatic32(temp[:32])
	forwardTransformStatic32(temp[32:])
	for i := 0; i < 32-1; i++ {
		input[i*2+0] = temp[i]
		input[i*2+1] = temp[i+32] + temp[i+32+1]
	}
	input[62], input[63] = temp[31], temp[63]
}

func forwardTransformStatic32(input []float64) {
	var temp [32]float64
	for i := 0; i < 16; i++ {
		x, y := input[i], input[31-i]
		temp[i] = x + y
		temp[i+16] = (x - y) / dct32[i]
	}
	forwardTransformStatic16(temp[:16])
	forwardTransformStatic16(temp[16:])
	for i := 0; i < 16-1; i++ {
		input[i*2+0] = temp[i]
		input[i*2+1] = temp[i+16] + temp[i+16+1]
	}

	input[30], input[31] = temp[15], temp[31]
}

func forwardTransformStatic16(input []float64) {
	var temp [16]float64
	for i := 0; i < 8; i++ {
		x, y := input[i], input[15-i]
		temp[i] = x + y
		temp[i+8] = (x - y) / dct16[i]
	}
	forwardTransformStatic8(temp[:8])
	forwardTransformStatic8(temp[8:])
	for i := 0; i < 8-1; i++ {
		input[i*2+0] = temp[i]
		input[i*2+1] = temp[i+8] + temp[i+8+1]
	}

	input[14], input[15] = temp[7], temp[15]
}

func forwardTransformStatic8(input []float64) {
	var temp [8]float64
	x0, y0 := input[0], input[7]
	x1, y1 := input[1], input[6]
	x2, y2 := input[2], input[5]
	x3, y3 := input[3], input[4]

	temp[0] = x0 + y0
	temp[1] = x1 + y1
	temp[2] = x2 + y2
	temp[3] = x3 + y3
	temp[4] = (x0 - y0) / dct8[0]
	temp[5] = (x1 - y1) / dct8[1]
	temp[6] = (x2 - y2) / dct8[2]
	temp[7] = (x3 - y3) / dct8[3]

	forwardTransformStatic4(temp[:4])
	forwardTransformStatic4(temp[4:])

	input[0] = temp[0]
	input[1] = temp[4] + temp[5]
	input[2] = temp[1]
	input[3] = temp[5] + temp[6]
	input[4] = temp[2]
	input[5] = temp[6] + temp[7]
	input[6] = temp[3]
	input[7] = temp[7]
}

func forwardTransformStatic4(input []float64) {
	var (
		t0, t1, t2, t3 float64
	)
	x0, y0 := input[0], input[3]
	x1, y1 := input[1], input[2]

	t0 = x0 + y0
	t1 = x1 + y1
	t2 = (x0 - y0) / dct4[0]
	t3 = (x1 - y1) / dct4[1]

	x, y := t0, t1
	t0 += t1
	t1 = (x - y) / dct2[0]

	x, y = t2, t3
	t2 += t3
	t3 = (x - y) / dct2[0]

	input[0] = t0
	input[1] = t2 + t3
	input[2] = t1
	input[3] = t3
}

func init() {
	// dct256
	for i := 0; i < 128; i++ {
		dct256[i] = (math.Cos((float64(i)+0.5)*math.Pi/256) * 2)
	}
	// dct128
	for i := 0; i < 64; i++ {
		dct128[i] = (math.Cos((float64(i)+0.5)*math.Pi/128) * 2)
	}
}

// Static DCT Tables
var (
	dct256 = [128]float64{}
	dct128 = [64]float64{}
	dct64  = [32]float64{
		(math.Cos((float64(0)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(1)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(2)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(3)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(4)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(5)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(6)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(7)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(8)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(9)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(10)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(11)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(12)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(13)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(14)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(15)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(16)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(17)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(18)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(19)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(20)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(21)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(22)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(23)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(24)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(25)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(26)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(27)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(28)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(29)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(30)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(31)+0.5)*math.Pi/64) * 2),
	}
	dct32 = [16]float64{
		(math.Cos((float64(0)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(1)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(2)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(3)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(4)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(5)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(6)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(7)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(8)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(9)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(10)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(11)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(12)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(13)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(14)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(15)+0.5)*math.Pi/32) * 2),
	}
	dct16 = [8]float64{
		(math.Cos((float64(0)+0.5)*math.Pi/16) * 2),
		(math.Cos((float64(1)+0.5)*math.Pi/16) * 2),
		(math.Cos((float64(2)+0.5)*math.Pi/16) * 2),
		(math.Cos((float64(3)+0.5)*math.Pi/16) * 2),
		(math.Cos((float64(4)+0.5)*math.Pi/16) * 2),
		(math.Cos((float64(5)+0.5)*math.Pi/16) * 2),
		(math.Cos((float64(6)+0.5)*math.Pi/16) * 2),
		(math.Cos((float64(7)+0.5)*math.Pi/16) * 2),
	}
	dct8 = [4]float64{
		(math.Cos((float64(0)+0.5)*math.Pi/8) * 2),
		(math.Cos((float64(1)+0.5)*math.Pi/8) * 2),
		(math.Cos((float64(2)+0.5)*math.Pi/8) * 2),
		(math.Cos((float64(3)+0.5)*math.Pi/8) * 2),
	}
	dct4 = [2]float64{
		(math.Cos((float64(0)+0.5)*math.Pi/4) * 2),
		(math.Cos((float64(1)+0.5)*math.Pi/4) * 2),
	}
	dct2 = [1]float64{
		(math.Cos((float64(0)+0.5)*math.Pi/2) * 2),
	}
)
