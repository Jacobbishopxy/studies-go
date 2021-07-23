package sample

import (
	"math/rand"
	"time"

	"pcbook/pb"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 返回一个新的样本 keyboard
func NewKeyboard() *pb.Keyboard {
	keyboard := &pb.Keyboard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}

	return keyboard
}

func randomBool() bool {
	return rand.Intn(2) == 1
}

func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}
}

func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

func randomStringFromSet(a ...string) string {
	n := len(a)
	if n == 0 {
		return ""
	}
	return a[rand.Intn(n)]
}

func randomCPUName(brand string) string {
	if brand == "Intel" {
		return randomStringFromSet(
			"Xeon E-2286M",
			"Core i9-9980HK",
			"Core i7-9750H",
			"Core i5-9400F",
			"Core i3-1005G1",
		)
	}

	return randomStringFromSet(
		"Ryzen 7 PRO 2700U",
		"Ryzen 5 PRO 3500U",
		"Ryzen 3 PRO 3200GE",
	)
}

func randomInt(min, max int) int {
	return min + rand.Int()%(max-min+1)
}

func randomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// 样本 CPU
func NewCPU() *pb.CPU {
	brand := randomCPUBrand()
	name := randomCPUName(brand)

	numberCores := randomInt(2, 8)
	numberThreads := randomInt(numberCores, 12)

	minGhz := randomFloat64(2.0, 3.5)
	maxGhz := randomFloat64(minGhz, 5.0)

	cpu := &pb.CPU{
		Brand:         brand,
		Name:          name,
		NumberCores:   uint32(numberCores),
		NumberThreads: uint32(numberThreads),
		MinGhz:        minGhz,
		MaxGhz:        maxGhz,
	}

	return cpu
}

func randomGPUBrand() string {
	return randomStringFromSet("Nvidia", "AMD")
}

func randomGPUName(brand string) string {
	if brand == "Nvidia" {
		return randomStringFromSet(
			"RTX 2060",
			"RTX 2070",
			"GTX 1660-Ti",
			"GTX 1070",
		)
	}

	return randomStringFromSet(
		"RX 580",
		"RX 590",
		"RX 5700-XT",
		"RX Vega-56",
	)
}

// 样本 GPU
func NewGPU() *pb.GPU {
	brand := randomGPUBrand()
	name := randomGPUName(brand)

	minGhz := randomFloat64(1.0, 1.5)
	maxGnz := randomFloat64(minGhz, 2.0)
	memGB := randomInt(2, 6)

	gpu := &pb.GPU{
		Brand:  brand,
		Name:   name,
		MinGhz: minGhz,
		MaxGhz: maxGnz,
		Memory: &pb.Memory{
			Value: uint64(memGB),
			Unit:  pb.Memory_GIGABYTE,
		},
	}

	return gpu
}

// 样本 RAM
func NewRAM() *pb.Memory {
	memGB := randomInt(4, 64)

	ram := &pb.Memory{
		Value: uint64(memGB),
		Unit:  pb.Memory_GIGABYTE,
	}

	return ram
}

// 样本 SSD
func NewSSD() *pb.Storage {
	memGB := randomInt(128, 1024)

	ssd := &pb.Storage{
		Driver: pb.Storage_SSD,
		Meomory: &pb.Memory{
			Value: uint64(memGB),
			Unit:  pb.Memory_GIGABYTE,
		},
	}

	return ssd
}

// 样本 HDD
func NewHDD() *pb.Storage {
	memTB := randomInt(1, 6)

	hdd := &pb.Storage{
		Driver: pb.Storage_HDD,
		Meomory: &pb.Memory{
			Value: uint64(memTB),
			Unit:  pb.Memory_TERABYTE,
		},
	}

	return hdd
}

func randomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func randomScreenResolution() *pb.Screen_Resolution {
	height := randomInt(1080, 4320)
	width := height * 16 / 9

	resolution := &pb.Screen_Resolution{
		Width:  uint32(width),
		Height: uint32(height),
	}
	return resolution
}

func randomScreenPanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLED
}

// 样本 Screen
func NewScreen() *pb.Screen {
	screen := &pb.Screen{
		SizeInch:   randomFloat32(13, 17),
		Resolution: randomScreenResolution(),
		Panel:      randomScreenPanel(),
		Multitouch: randomBool(),
	}

	return screen
}

func randomID() string {
	return uuid.New().String()
}

func randomLaptopBrand() string {
	return randomStringFromSet("Apple", "Dell", "Windows")
}

func randomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFromSet("Macbook Air", "Macbook Pro")
	case "Dell":
		return randomStringFromSet("Latitude", "Vostro", "XPS", "Alienware")
	default:
		return randomStringFromSet("Surface Book", "Surface Pro", "Surface Laptop")
	}
}

func NewLaptop() *pb.Laptop {
	brand := randomLaptopBrand()
	name := randomLaptopName(brand)

	laptop := &pb.Laptop{
		Id:       randomID(),
		Brand:    brand,
		Name:     name,
		Cpu:      NewCPU(),
		Ram:      NewRAM(),
		Gpus:     []*pb.GPU{NewGPU()},
		Storages: []*pb.Storage{NewSSD(), NewHDD()},
		Screen:   NewScreen(),
		Keyboard: NewKeyboard(),
		Weight: &pb.Laptop_WeightKg{
			WeightKg: randomFloat64(1500, 3500),
		},
		PriceUsed:   randomFloat64(1500, 3500),
		ReleaseYear: uint32(randomInt(2015, 2021)),
		UpdatedAt:   timestamppb.Now(),
	}

	return laptop
}
