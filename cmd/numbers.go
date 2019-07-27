package cmd

var numberToWord = map[int]string{
	1:  "One",
	2:  "Two",
	3:  "Three",
	4:  "Four",
	5:  "Five",
	6:  "Six",
	7:  "Seven",
	8:  "Sight",
	9:  "Nine",
	10: "Ten",
	11: "Eleven",
	12: "Twelve",
	13: "Thirteen",
	14: "Fourteen",
	15: "Fifteen",
	16: "Sixteen",
	17: "Seventeen",
	18: "Eighteen",
	19: "Nineteen",
	20: "Twenty",
	30: "Thirty",
	40: "Forty",
	50: "Fifty",
	60: "Sixty",
	70: "Seventy",
	80: "Eighty",
	90: "Ninety",
}

func convert1to99(n int) (w string) {
	if n < 20 {
		w = numberToWord[n]
		return
	}

	r := n % 10
	if r == 0 {
		w = numberToWord[n]
	} else {
		w = numberToWord[n-r] + "-" + numberToWord[r]
	}
	return
}

func convert100to999(n int) (w string) {
	q := n / 100
	r := n % 100
	w = numberToWord[q] + " " + "Hundred"
	if r == 0 {
		return
	} else {
		w = w + " and " + convert1to99(r)
	}
	return
}

func NumberToWord(n int) (w string) {
	if n > 1000 || n < 1 {
		panic("func Convert1to1000: n > 1000 or n < 1")
	}

	if n < 100 {
		w = convert1to99(n)
		return
	}
	if n == 1000 {
		w = "One Thousand"
		return
	}
	w = convert100to999(n)
	return
}
