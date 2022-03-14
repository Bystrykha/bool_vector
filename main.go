package main

import (
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"time"
)

type BF struct {
	variablesNumb, sliceLen int      //количество переменных / длина среза, в котором содержаться биты функции
	functionValue           []uint32 // срез с битами (используется именно срез, поскольку из него можно удалять/добавлять элементы)
}

func (b BF) newBFArgs2(variablesNumb int, sliceLen int, functionValue []uint32) BF { // конструктор по умолчанию
	b.variablesNumb = variablesNumb
	b.sliceLen = sliceLen
	b.functionValue = functionValue
	return b
}

func (b BF) newBF() BF { // конструктор по умолчанию
	b.variablesNumb = 0
	b.sliceLen = 0
	return b
}

func fillFunc(len int, function []uint32, value int) []uint32 { //заполнение функции только 0 или 1
	for i := 0; i < len; i++ {
		function = append(function, 0)
	}
	if value == 1 {
		for i := 0; i < len; i++ {
			function[i] -= 1
		}
	}
	return function
}

func randFunc(len int, function []uint32) []uint32 { // заполнение функции случайными значениями
	time.Sleep(10) // поскольку код выполняется быстро, при создании 2-х объектов может получиться одинаковый результат, чтобы этого избежать используется задрержка
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len; i++ {
		function = append(function, uint32(rand.Intn(4294967296)))
	}
	return function
}

func (b BF) newBFArgs(variablesNumb int, value int) (BF, error) { // конструктор с аргументами
	b.variablesNumb = variablesNumb
	b.sliceLen = int((1 << b.variablesNumb) >> 5)
	if b.sliceLen == 0 {
		b.sliceLen = 1
	}
	if value == 0 || value == 1 {
		b.functionValue = fillFunc(b.sliceLen, b.functionValue, value)
	} else if value == 2 {
		b.functionValue = randFunc(b.sliceLen, b.functionValue)
	} else {
		return b, errors.New("Incorrect input data.")
	}
	if variablesNumb < 5 {
		b.functionValue[0] %= uint32(1 << (1 << variablesNumb))
	}
	return b, nil
}

func (b BF) copyBF(pattern BF) BF { // конструктор копирования
	b.variablesNumb = pattern.variablesNumb
	b.sliceLen = pattern.sliceLen
	b.functionValue = nil
	b.functionValue = append(b.functionValue, pattern.functionValue...)
	return b
}

func reverse(s string) string { // разворот строки
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func (b BF) stringToBF(s string) BF { // конструктор по строке (если длина строки != степени двойки, считаем, что недостающие символы - нули)
	if s == "" {
		return b.newBF()
	} else {
		for i := 1; true; i++ {
			if int(1<<i) >= len(s) {
				b.variablesNumb = i
				break
			}
		}
		for i, v := range s {
			shift := i % 32
			if shift == 0 {
				b.functionValue = append(b.functionValue, 0)
			}
			if v == 49 {
				b.functionValue[len(b.functionValue)-1] |= 1 << shift
			}
		}
		b.sliceLen = len(b.functionValue)
		return b
	}
}

func (b BF) getWeight() int { // подсчет веса вектора
	weight := 0
	for _, v := range b.functionValue {
		for ; v > 0; v = v & (v - 1) {
			weight += 1
		}
	}
	return weight
}

func (b BF) compare(b2 BF) bool { // сравнение векторов (только ==)
	if b.sliceLen != b2.sliceLen {
		return false
	}
	for i := range b.functionValue {
		if b.functionValue[i] != b2.functionValue[i] {
			return false
		}
	}
	return true
}

func (b BF) xor(b2 BF) BF { // сложение по модулю 2
	resVector := b2.functionValue
	minLen := b.sliceLen
	variablesNumb := b2.variablesNumb
	if b2.sliceLen < minLen {
		minLen = b2.sliceLen
		resVector = b.functionValue
		variablesNumb = b.variablesNumb
	}
	for i := 0; i < minLen; i++ {
		resVector[i] = b.functionValue[i] ^ b2.functionValue[i]
	}
	var res BF
	res = res.newBFArgs2(variablesNumb, len(resVector), resVector)
	return res
}

func (b BF) logMul(b2 BF) BF { // логическое умножение
	resVector := b.functionValue
	minLen := b.sliceLen
	variablesNumb := b.variablesNumb
	if b2.sliceLen < minLen {
		minLen = b2.sliceLen
		resVector = b2.functionValue
		variablesNumb = b2.variablesNumb
	}
	for i := 0; i < minLen; i++ {
		resVector[i] = b.functionValue[i] & b2.functionValue[i]
	}
	var res BF
	res = res.newBFArgs2(variablesNumb, len(resVector), resVector)
	return res
}

func (b BF) logAdd(b2 BF) BF { //логическое сложение
	resVector := b2.functionValue
	minLen := b.sliceLen
	variablesNumb := b2.variablesNumb
	if b2.sliceLen < minLen {
		minLen = b2.sliceLen
		resVector = b.functionValue
		variablesNumb = b.variablesNumb
	}
	for i := 0; i < minLen; i++ {
		resVector[i] = b.functionValue[i] | b2.functionValue[i]
	}
	var res BF
	res = res.newBFArgs2(variablesNumb, len(resVector), resVector)
	return res
}

func (b BF) leftShift(shift int) BF { // сдвиг влево
	if shift <= 32 {
		shiftBuffer1 := uint32(0)
		shiftBuffer2 := uint32(0)
		for i := range b.functionValue {
			shiftBuffer1 = b.functionValue[i]
			b.functionValue[i] <<= shift
			b.functionValue[i] |= shiftBuffer2
			shiftBuffer2 = shiftBuffer1 >> (32 - shift)
		}
	} else {
		a := shift / 32
		for i := 0; i < b.sliceLen-a; i++ {
			b.functionValue[i+a] = b.functionValue[i]
		}
		for i := 0; i < a; i++ {
			b.functionValue[i] = uint32(0)
		}
	}
	return b
}

func (b BF) rightShift(shift int) BF { // сдвиг вправо
	if shift <= 32 {
		shiftBuffer1 := uint32(0)
		shiftBuffer2 := uint32(0)
		for i := b.sliceLen - 1; i >= 0; i-- {
			shiftBuffer1 = b.functionValue[i]
			b.functionValue[i] >>= shift
			b.functionValue[i] |= shiftBuffer2
			shiftBuffer2 = shiftBuffer1 << (32 - shift)
		}
	} else {
		a := shift / 32
		for i := a; i < b.sliceLen; i++ {
			b.functionValue[i] = b.functionValue[i-a]
		}
		for i := b.sliceLen - 1; i < b.sliceLen-a; i++ { //мб придктся исправить на i < b.sliceLen-a-1
			b.functionValue[i] = uint32(0)
		}
	}
	return b
}

func (b BF) outVector() { // красивый вывод вектора
	var res string
	for i := range b.functionValue {
		res += reverse(fmt.Sprintf("%032b", b.functionValue[i]))
	}
	if b.variablesNumb < 5 {
		res = res[:(1 << b.variablesNumb)]
	}
	fmt.Println(res)
}

func (b BF) indexesOne() []int {
	var indexesOne []int
	k := 0
	for _, v := range b.functionValue {
		var a uint32 = 1
		for j := 0; j < 32; j++ {
			if v&a != 0 {
				indexesOne = append(indexesOne, k+j)
			}
			a <<= 1
		}
		k += 32
	}
	return indexesOne
}

func (b BF) getANF() string {
	g := b.getMobius()
	indexesOne := g.indexesOne()
	var s string
	if len(indexesOne) > 0 {
		a := 0
		if indexesOne[0] == 0 {
			s += "1+"
			a = 1
		}
		for i := a; i < len(indexesOne); i++ {
			vector := 1 << g.variablesNumb
			for j := 0; vector > 0; j++ {
				vector >>= 1
				if indexesOne[i]&vector > 0 {
					s = s + "x" + strconv.FormatInt(int64(j+1), 10)
				}
			}
			s += "+"
		}
		s = s[0 : len(s)-1]
	} else if len(indexesOne) == 0 {
		s = "0"
	}
	return s
}

func getMulVector(shift int, len int) BF {
	var vector BF
	vector, _ = vector.newBFArgs(len, 0)
	if shift == 0 {
		a := uint32(2863311530)
		for i := range vector.functionValue {
			vector.functionValue[i] = a
		}
	} else if shift == 1 {
		a := uint32(3435973836)
		for i := range vector.functionValue {
			vector.functionValue[i] = a
		}
	} else if shift == 2 {
		a := uint32(4042322160)
		for i := range vector.functionValue {
			vector.functionValue[i] = a
		}
	} else if shift == 3 {
		a := uint32(4278255360)
		for i := range vector.functionValue {
			vector.functionValue[i] = a
		}
	} else if shift == 4 {
		a := uint32(4294901760)
		for i := range vector.functionValue {
			vector.functionValue[i] = a
		}
	}
	return vector
}

func (b BF) getMobius() BF {
	var g, s, m, stepVector BF
	g = b
	for i := 0; i < b.variablesNumb; i++ {
		if i < 5 {
			s = s.copyBF(g)
			s = s.leftShift(1 << i)
			stepVector = getMulVector(i, b.variablesNumb)
			m = s.logMul(stepVector)
			g = g.xor(m)
		} else {
			k := 1 << (i - 5)
			for j := 0; j < len(g.functionValue); {
				for l := 0; l < k; l += 1 {
					g.functionValue[j+k] = g.functionValue[j] ^ g.functionValue[j+k]
					j += 1
				}
				j += k
			}
		}
	}
	return g
}

func (b BF) getDegree() int {
	g := b.getMobius()
	indexesOne := g.indexesOne()
	functionDegree, c := 0, g.variablesNumb
	for _, v := range indexesOne {
		monomWeight := 0
		for c >= 1 {
			if v%(1<<c) >= (1 << (c - 1)) {
				monomWeight += 1
			}
			c -= 1
		}
		if monomWeight > functionDegree {
			functionDegree = monomWeight
		}
	}
	return functionDegree
}

func main() {
	var b BF
	b = b.stringToBF("00011110")
	fmt.Println(b.getANF())
	fmt.Println(b.getDegree())
}
