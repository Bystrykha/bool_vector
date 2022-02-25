package main

import (
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"time"
)

type BF struct {
	variablesNumb, sliceLen int		//количество переменных / длина среза, в котором содержаться биты функции
	functionValue []uint32	// срез с битами (используется именно срез, поскольку из него можно удалять/добавлять элементы)
}

func (b BF) newBF() BF{		// конструктор по умолчанию
	b.variablesNumb = 0
	b.sliceLen = 0
	return b
}

func fillFunc(len int, function []uint32, value int) []uint32{	//заполнение функции только 0 или 1
	for i:=0; i < len; i++{
		function = append(function, 0)
	}
	if value == 1{
		for i:=0; i < len; i++{
			function[i] -= 1
		}
	}
	return function
}

func randFunc(len int, function []uint32) []uint32{		// заполнение функции случайными значениями
	time.Sleep(10)		// поскольку код выполняется быстро, при создании 2-х объектов может получиться одинаковый результат, чтобы этого избежать используется задрержка
	rand.Seed(time.Now().UnixNano())
	for i:=0; i < len; i++ {
		function = append(function, uint32(rand.Intn(4294967296)))
	}
	return function
}

func (b BF) newBFArgs (variablesNumb int, value int)(BF, error){	// конструктор с аргументами
	b.variablesNumb = variablesNumb
	b.sliceLen = int((1 << b.variablesNumb) >> 5)
	if b.sliceLen == 0{
		b.sliceLen = 1
	}
	if value == 0 || value == 1 {
		b.functionValue = fillFunc(b.sliceLen, b.functionValue, value)
	}else if value == 2{
			b.functionValue = randFunc(b.sliceLen, b.functionValue)
	}else {
		return b, errors.New("Incorrect input data.")
	}
	if variablesNumb < 5{
		b.functionValue[0] %= uint32(1 << (1 << variablesNumb))
	}
	return b, nil
}

func (b BF)copyBF(pattern BF)BF{	// конструктор копирования
	b.variablesNumb = pattern.variablesNumb
	b.sliceLen = pattern.sliceLen
	b.functionValue = pattern.functionValue
	return b
}

func reverse(s string) string {		// разворот строки
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func (b BF) stringToBF(s string) BF {	 // конструктор по строке (если длина строки != степени двойки, считаем, что недостающие символы - нули)
	if s == ""{
		return b.newBF()
	}else {
		s = reverse(s)
		for i:=1; true; i++{
			if int(1 << i) >= len(s){
				b.variablesNumb = i
				break
			}
		}
		for i, v := range s{
			shift := i % 32
			if shift == 0 {
				b.functionValue = append(b.functionValue, 0)
			}
			if v == 49{
				b.functionValue[len(b.functionValue)-1] |= 1 << shift
			}
		}
		b.sliceLen = len(b.functionValue)
		return b
	}
}

func (b BF) getWeight() int{		// подсчет веса вектора
	weight := 0
	for _, v := range b.functionValue{
		for ;v > 0; v = v & (v - 1){
			weight += 1
		}
	}
	return weight
}

func (b BF)compare(b2 BF)bool{		// сравнение векторов (только ==)
	if b.sliceLen != b2.sliceLen{
		return false
	}
	for i := range b.functionValue{
		if b.functionValue[i] != b2.functionValue[i]{
			return false
		}
	}
	return true
}

func (b BF) xor(b2 BF) []uint32{	// сложение по модулю 2
	resVector := b2.functionValue
	minLen := b.sliceLen
	if b2.sliceLen < minLen{
		minLen = b2.sliceLen
		resVector = b.functionValue
	}
	for i := 0; i < minLen; i++{
		resVector[i] = b.functionValue[i] ^ b2.functionValue[i]
	}
	return resVector
}

func (b BF) logMul(b2 BF) []uint32{ // логическое умножение
	resVector := b.functionValue
	minLen := b.sliceLen
	if b2.sliceLen < minLen{
		minLen = b2.sliceLen
		resVector = b2.functionValue
	}
	for i := 0; i < minLen; i++{
		resVector[i] = b.functionValue[i] & b2.functionValue[i]
	}
	return resVector
}

func (b BF) logAdd(b2 BF) []uint32{	//логическое сложение
	resVector := b2.functionValue
	minLen := b.sliceLen
	if b2.sliceLen < minLen{
		minLen = b2.sliceLen
		resVector = b.functionValue
	}
	for i := 0; i < minLen; i++{
		resVector[i] = b.functionValue[i] | b2.functionValue[i]
	}
	return resVector
}

func (b BF) leftShift(shift int) BF {		// сдвиг влево
	shiftBuffer1 := uint32(0)
	shiftBuffer2 := uint32(0)
	for i := range b.functionValue{
		shiftBuffer1 = b.functionValue[i]
		b.functionValue[i] <<= shift
		b.functionValue[i] |= shiftBuffer2
		shiftBuffer2 = shiftBuffer1 >> (32-shift)
	}
	return b
}

func (b BF) rightShift(shift int) BF {		// сдвиг вправо
	shiftBuffer1 := uint32(0)
	shiftBuffer2 := uint32(0)
	for i := b.sliceLen-1; i >= 0; i--{
		shiftBuffer1 = b.functionValue[i]
		b.functionValue[i] >>= shift
		b.functionValue[i] |= shiftBuffer2
		shiftBuffer2 = shiftBuffer1 << (32-shift)
	}
	return b
}

func (b BF) outVector(){		// красивый вывод вектора
	var res string
	for i := range b.functionValue{
		res += reverse(fmt.Sprintf("%032b", b.functionValue[i]))
	}
	if b.variablesNumb < 5{
		res = res[:(1 << b.variablesNumb)]
	}
	fmt.Println(res)
}

func main(){
	var b BF
	b, _ = b.newBFArgs(6, 2)
	fmt.Printf("%032b \n", b.functionValue)
	b.outVector()
	b = b.leftShift(3)
	b.outVector()
}