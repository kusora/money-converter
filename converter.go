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

	UNIT_FEN  = "分"
	UNIT_JIAO = "角"
	UNIT_YUAN = "元"

	POSTFIX = "整"
	PREFIX  = "￥"

	POS_UNIT_YUAN = 1
	POS_UNIT_WAN  = 5
	POS_UINT_YI   = 9
	MAX_POS       = 18
)

var unit_mapper = map[int]string{
	1: "",
	2: UNIT_SHI,
	3: UNIT_BAI,
	4: UNIT_QIAN,
	5: UNIT_WAN,
	9: UNIT_YI,
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

func Arab2Chinsese(num int64) (string, error) {
	// 最大值设为1000万亿
	numStr := strconv.FormatInt(num, 10)

	if len(numStr) > MAX_POS {
		return "", errors.New("too large number")
	}

	result := ""
	if len(numStr) < 3 {
		fenPart, _, err := parseUnitYuan(numStr)
		if err != nil {
			return "", err
		}
		result = fenPart
	} else {
		yuanPart, _, withZeroTail, err := parse(numStr[:len(numStr)-2])
		if err != nil {
			return "", err
		}
		fenPart, fenZeroHead, err := parseUnitYuan(numStr[len(numStr)-2:])
		if err != nil {
			return "", err
		}
		result = yuanPart + UNIT_YUAN
		if fenPart != "" && (withZeroTail == true || fenZeroHead == true) {
			result = result + num_mapper[0]
		}
		result = result + fenPart
	}
	result = PREFIX + result
	if numStr[len(numStr)-1] == '0' {
		result = result + POSTFIX
	}

	return result, nil

}

func getUnit(pos int) (string, error) {
	if pos > MAX_POS {
		return "", errors.New("too large amount")
	}

	value, ok := unit_mapper[pos]
	if !ok {
		if pos > POS_UINT_YI {
			return getUnit(pos - POS_UINT_YI + 1)
		}
		if pos > POS_UNIT_WAN {
			return getUnit(pos - POS_UNIT_WAN + 1)
		}
	}
	return value, nil
}

func parseUnitYuan(value string) (string, bool, error) {
	result := ""
	withZeroHead := false
	if len(value) == 0 {
		return "", true, nil
	}

	if len(value) == 1 {
		return num_mapper[int(value[0]-'0')] + UNIT_FEN, withZeroHead, nil
	}

	if value[0] != '0' {
		result = num_mapper[int(value[0]-'0')] + UNIT_JIAO
	} else {
		withZeroHead = true
	}

	if value[1] != '0' {
		result = result + num_mapper[int(value[1]-'0')] + UNIT_FEN
	}
	return result, withZeroHead, nil
}

func parse(input string) (result string, withZeroHead bool, withZeroTail bool, err error) {

	value, err := trimHead(input)
	if err != nil {
		return "", false, false, err
	}
	if len(value) == 0 {
		return "", true, true, nil
	}
	result = ""
	withZeroHead = value[0] == '0'
	withZeroTail = value[len(value)-1] == '0'
	// 判断是不是大于yi,
	if len(value) > POS_UINT_YI {
		yiValue, _, yiZeroTail, newerr := parse(value[:len(value)-POS_UINT_YI+1])
		if newerr != nil {
			err = newerr
			return
		}
		result = yiValue + unit_mapper[POS_UINT_YI]
		yuanValue, yuanZeroHead, _, newerr := parse(value[len(value)-POS_UINT_YI+1:])
		if newerr != nil {
			err = newerr
			return
		}
		if yuanValue != "" && (yiZeroTail == true || yuanZeroHead == true) {
			result = result + num_mapper[0]
		}
		result = result + yuanValue
		return
	}
	if len(value) > POS_UNIT_WAN {
		wanValue, _, wanZeroTail, newerr := parse(value[:len(value)-POS_UNIT_WAN+1])
		if newerr != nil {
			err = newerr
			return
		}
		result = wanValue + unit_mapper[POS_UNIT_WAN]
		yuanValue, yuanZeroHead, _, newerr := parse(value[len(value)-POS_UNIT_WAN+1:])
		if newerr != nil {
			err = newerr
			return
		}
		if yuanValue != "" && (wanZeroTail == true || yuanZeroHead == true) {
			result = result + num_mapper[0]
		}
		result = result + yuanValue
		return
	}
	if len(value) > 0 {
		result = ""
		zeroBefore := true
		for i := 1; i <= len(value); i++ {
			if value[len(value)-i] != '0' {
				result = num_mapper[int(value[len(value)-i]-'0')] + unit_mapper[i] + result
				zeroBefore = false
			} else {
				if zeroBefore == false {
					result = num_mapper[0] + result
					zeroBefore = true
				}
			}
		}
		return
	}
	return
}

func trimHead(input string) (string, error) {
	for i := 0; i < len(input); i++ {
		if input[i] != '0' {
			return input[i:], nil
		}
	}
	return "", nil
}
