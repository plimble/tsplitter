package tsplitter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplit(t *testing.T) {
	txt := `โดราเอมอน (ญี่ปุ่น: ドラえもん Dora'emon โดะระเอะมง ?) หรือ โดเรมอน เป็น การ์ตูนญี่ปุ่น แต่งโดย ฟุจิโกะ ฟุจิโอะ เรื่องราวของหุ่นยนต์แมวหูด้วน ชื่อ โดราเอมอน โดยฟุจิโกะ ฟุจิโอะ ได้กล่าวว่าโดราเอมอนเกิดวันที่ 3 กันยายน ค.ศ. 2112 (พ.ศ. 2655) มาจากอนาคตเพื่อกลับมาช่วยเหลือโนบิตะ เด็กประถมจอมขี้เกียจด้วยของวิเศษจากอนาคต โดราเอมอนเริ่มตีพิมพ์ครั้งแรกในญี่ปุ่นเมื่อเดือนมกราคม พ.ศ. 2513 โดยสำนักพิมพ์โชงะกุกัง[3][4] โดยมีจำนวนตอนทั้งหมด 1,344 ตอน[5] ต่อมาในวันที่ 11 มิถุนายน พ.ศ. 2540 โดราเอมอนได้รับรางวัลเทะซุกะ โอซามุ ครั้งที่ 1 ในสาขาการ์ตูนดีเด่น[6] อีกทั้งยังได้รับเลือกจากนิตยสารไทม์เอเชีย ให้เป็นหนึ่งในวีรบุรุษของทวีปเอเชีย จากประเทศญี่ปุ่น[7] จากนั้นในวันที่ 19 มีนาคม พ.ศ. 2551 โดราเอมอนก็ได้รับเลือกให้เป็นทูตสันถวไมตรีเพื่อการประชาสัมพันธ์วัฒนธรรมของประเทศญี่ปุ่น[8] นอกจากนี้บริษัทบันได ผู้ผลิตและจำหน่ายสินค้าการ์ตูนที่มีชื่อเสียงของญี่ปุ่น ยังได้ผลิตหุ่นยนต์โดราเอมอนของจริงขึ้นมาในชื่อว่า "My โดราเอมอน" โดยออกวางจำหน่ายครั้งแรกในวันที่ 3 กันยายน พ.ศ. 2552[9]

	ในประเทศไทย โดราเอมอนฉบับหนังสือการ์ตูนมีการตีพิมพ์โดยหลายสำนักพิมพ์ในช่วงก่อนที่จะมีลิขสิทธิ์การ์ตูน[10][11] แต่ปัจจุบัน สำนักพิมพ์เนชั่น เอ็ดดูเทนเมนท์ เป็นผู้ได้รับลิขสิทธิ์ในการจัดพิมพ์แต่เพียงผู้เดียว ส่วนฉบับอะนิเมะ ออกอากาศทางช่อง 9 อ.ส.ม.ท. หรือโมเดิร์นไนน์ทีวี ในปัจจุบัน และวางจำหน่ายในรูปแบบวีซีดี-ดีวีดี ลิขสิทธิ์โดยบริษัทโรส วิดีโอ[12] แพทช์ แพทช์ าแ ่อา ฐาน าฐาน าฐานฐาน ฮฮฮฮฮฮฮ ันักศึกษาใจใหญ่ ส้ม ใส้ม ักาฎ`

	dict := NewFileDict("dictionary.txt")

	s := Split(dict, txt)
	fmt.Println(s.Size())
	fmt.Println("")
	fmt.Println("All", len(s.All()))
	fmt.Println(s.All())
	fmt.Println("")
	fmt.Println("Unknown", len(s.Unknown()))
	fmt.Println(s.Unknown())
	assert.Len(t, s.All(), 262)
	assert.Len(t, s.Unknown(), 4)
}

func TestIsEnglish(t *testing.T) {
	assert.True(t, isEnglish('a'))
	assert.True(t, isEnglish('Z'))
}

func TestIsDigit(t *testing.T) {
	assert.True(t, isDigit('1'))
	assert.True(t, isDigit('๙'))
}

func TestIsSpecialChar(t *testing.T) {
	assert.True(t, isSpecialChar('~'))
	assert.True(t, isSpecialChar('ๆ'))
	assert.True(t, isSpecialChar('ฯ'))
	assert.True(t, isSpecialChar('“'))
	assert.True(t, isSpecialChar('”'))
	assert.True(t, isSpecialChar(','))
}

func TestIsEnding(t *testing.T) {
	assert.True(t, isEnding('ๆ'))
	assert.True(t, isEnding('ฯ'))
	assert.False(t, isEnding('?'))
}
