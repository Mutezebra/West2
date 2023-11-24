package myutils

func StringToUint(str string) uint {
	var num uint
	for _, v := range str {
		num = num*10 + uint(v-'0')
	}
	return num
}

func UintToString(num uint) string {
	var str string
	for num > 0 {
		str = string(num%10+'0') + str
		num /= 10
	}
	return str
}
