package converter

import (
	"testing"
)

func TestConvert(t *testing.T) {
	check(100000000000000000, "￥壹仟万亿元整", t)
	check(100000000000010001, "￥壹仟万亿零壹佰元零壹分", t)
	check(111111111111111111, "￥壹仟壹佰壹拾壹万壹仟壹佰壹拾壹亿壹仟壹佰壹拾壹万壹仟壹佰壹拾壹元壹角壹分", t)
	check(123456789012345678, "￥壹仟贰佰叁拾肆万伍仟陆佰柒拾捌亿玖仟零壹拾贰万叁仟肆佰伍拾陆元柒角捌分", t)
	check(101010101010101001, "￥壹仟零壹拾万零壹仟零壹拾亿零壹仟零壹拾万零壹仟零壹拾元零壹分", t)
	check(1, "￥壹分", t)
	check(10, "￥壹角整", t)
	check(100, "￥壹元整", t)
	check(101, "￥壹元零壹分", t)
	check(1000, "￥壹拾元整", t)
	check(1100, "￥壹拾壹元整", t)
	check(10100, "￥壹佰零壹元整", t)
	check(11000, "￥壹佰壹拾元整", t)
	check(11100, "￥壹佰壹拾壹元整", t)
	check(101100, "￥壹仟零壹拾壹元整", t)
	check(168032, "￥壹仟陆佰捌拾元零叁角贰分", t)
	check(10700053, "￥壹拾万零柒仟元零伍角叁分", t)

}

func check(num int64, expect string, t *testing.T) {
	actual, err := Arab2Chinsese(num)
	if err != nil {

	}
	if actual != expect {
		t.Error("not match", actual, expect)
	}
}
