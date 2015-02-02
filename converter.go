package converter

import (
	"errors"
	"strconv"
)

const (
	UNIT_SHI  = "拾"
	UNIT_BAI  = "佰"
	UNIT_QIAN = "仟"

	UNIT_WAN  = "万"
	UNIT_YI   = "亿"
	UNIT_YUAN = "元"

	UNIT_FEN  = "分"
	UNIT_JIAO = "角"

	POSTFIX = "整"
	PREFIX  = "￥"

	BASE_PART_LENTH = 4
	MAX_POS         = 18
)

var low_unit_mapper = map[int]string{
	1: "",
	2: UNIT_SHI,
	3: UNIT_BAI,
	4: UNIT_QIAN,
}

var high_unit_mapper = map[int]string{
	0: UNIT_YUAN,
	1: UNIT_WAN,
	2: UNIT_YI,
	3: UNIT_WAN, // 万亿
}

var num_mapper = map[int]string{
	1: "壹",
	2: "贰",
	3: "叁",
	4: "肆",
	5: "伍",
	6: "陆",
	7: "柒",
	8: "捌",
	9: "玖",
	0: "零",
}

func Arab2Chinsese2(num int64) (string, error) {
	// 最大值设为1000万亿
	numStr := strconv.FormatInt(num, 10)

	if len(numStr) > MAX_POS {
		return "", errors.New("too large number")
	}

	pieces := make([]string, 0, 4)
	whole := numStr

	for {
		if len(whole) > 2 {
			length := 4
			if (len(whole)-2)%BASE_PART_LENTH != 0 {
				length = (len(whole) - 2) % BASE_PART_LENTH
			}
			pieces = append(pieces, whole[0:length])
			whole = whole[length:]
		} else {
			pieces = append(pieces, whole)
			break
		}
	}

	var result, before string
	for index, piece := range pieces {
		if index != len(pieces)-1 {
			currentValue := baseParse(piece)
			if currentValue != "" && withZero(before, piece) {
				result = result + num_mapper[0] + currentValue + high_unit_mapper[len(pieces)-index-2]
			} else if currentValue != "" || index != len(pieces)-3 {
				result = result + currentValue + high_unit_mapper[len(pieces)-index-2]
			}
		} else {
			currentValue := lowParse(piece)
			if currentValue != "" && withZero(before, piece) {
				result = result + num_mapper[0] + currentValue
			} else {
				result = result + currentValue
			}
			break
		}
		before = piece
	}
	result = PREFIX + result

	if result == "" {
		result = num_mapper[0] + UNIT_YUAN + POSTFIX
	} else if numStr[len(numStr)-1] == '0' {
		result = result + POSTFIX
	}

	return result, nil

}

func lowParse(value string) string {
	switch len(value) {
	case 0:
		return ""
	case 1:
		return num_mapper[int(value[0]-'0')] + UNIT_FEN
	default:
		result := ""
		if value[0] != '0' {
			result = num_mapper[int(value[0]-'0')] + UNIT_JIAO
		}

		if value[1] != '0' {
			result = result + num_mapper[int(value[1]-'0')] + UNIT_FEN
		}
		return result
	}
}

// 转换小于10000的数字， 属于基础函数
func baseParse(input string) string {
	result, zero := "", ""
	lastValue := '0'
	for index, value := range input {
		if value != '0' {
			result = result + zero + num_mapper[int(value-'0')] + low_unit_mapper[len(input)-index]
			lastValue, zero = value, ""
		} else {
			if lastValue != '0' {
				zero = num_mapper[0]
			}
			lastValue = value
		}
	}
	return result
}

func withZero(before string, after string) bool {
	if len(before) == 0 || len(after) == 0 {
		return false
	} else if before[len(before)-1] == '0' || after[0] == '0' {
		return true
	}
	return false
}
