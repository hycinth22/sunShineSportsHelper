package utility

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

var UA = []string{
	"Dalvik/2.1.0 (Linux; U; Android 7.0; GT-I9152P Build/JLS36C)",
	"UCWEB7.0.2.37/28/999",
	"Openwave/ UCWEB7.0.2.37/28/999",
	"Mozilla/4.0 (compatible; MSIE 6.0; ) Opera/UCWEB7.0.2.37/28/999",
	"Mozilla/5.0 (BlackBerry; U; BlackBerry 9800; en) AppleWebKit/534.1+ (KHTML, like Gecko) Version/6.0.0.337 Mobile Safari/534.1+",
	"Mozilla/5.0 (Linux; U; Android 3.0; en-us; Xoom Build/HRI39) AppleWebKit/534.13 (KHTML, like Gecko) Version/4.0 Safari/534.13",
	"Opera/9.80 (Android 2.3.4; Linux; Opera Mobi/build-1107180945; U; en-GB) Presto/2.8.149 Version/11.10",
	"MQQBrowser/26 Mozilla/5.0 (Linux; U; Android 2.3.7; zh-cn; MB200 Build/GRJ22; CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
	"Mozilla/5.0 (Linux; U; Android 2.3.7; en-us; Nexus One Build/FRF91) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
	"Mozilla/5.0 (Linux; U; Android 4.0.3; zh-cn; M032 Build/IML74K) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/4.1 Mobile Safari/533.1",
	"Mozilla/5.0 (Linux; U; Android 4.0.3; zh-cn; M032 Build/IML74K) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
	"Opera/9.80 (Android 4.0.3; Linux; Opera Mobi/ADR-1210241554) Presto/2.11.355 Version/12.10",
	"Nokia5320/04.13 (SymbianOS/9.3; U; Series60/3.2 Mozilla/5.0; Profile/MIDP-2.1 Configuration/CLDC-1.1 ) AppleWebKit/413 (KHTML, like Gecko) Safari/413",
	"Nokia5320(19.01)/SymbianOS/9.1 Series60/3.0",
	"Dalvik/2.3.0 (Linux; U; Android 7.0; Redmi Note 2 Build/LRX22G)",
	"Dalvik/2.0.3 (Linux; U; Android 6.0; HUAWEI G6-C00 Build/HuaweiG6-C00)",
	"Dalvik/1.9.0 (Linux; U; Android 5.0; OPPO R7 Build/KTU84P)",
	"Dalvik/1.9.0 (Linux; U; Android 5.0; OPPO R7 Build/KTU84P)",
	"Dalvik/1.2.0+(Linux;+U;+Android+2.3.5;+LX6200B+Build/MocorDroid2.3.5)",
	"Dalvik/1.2.0+(Linux;+U;+Android+2.3.5;+QX333+Build/MocorDroid2.3.5)",
	"Dalvik/1.2.0+(Linux;+U;+Android+2.3.5;+android+Build/MocorDroid2.3.5)",
	"Dalvik/1.2.0+(Linux;+U;+Android+2.3.6;+P9220+Build/MocorDroid2.3.6)",
	"Dalvik/1.4.0+(Linux;+U;+Android+2.3.4;+MT15i+Build/4.0.2.A.0.62)",
	"Dalvik/1.4.0+(Linux;+U;+Android+2.3.5;+999+Build/GRJ90)",
	"Dalvik/1.4.0+(Linux;+U;+Android+2.3.5;+ChangHong+V7+Build/GRJ90)",
	"Dalvik/1.4.0+(Linux;+U;+Android+2.3.5;+DESAY+TS808+Build/MocorDroid2.3.5)",
	"Dalvik/1.4.0+(Linux;+U;+Android+2.3.5;+LT+S600D+Build/GRJ90)",
	"Dalvik/1.4.0+(Linux;+U;+Android+2.3.5;+T100+Build/MocorDroid2.3.5)",
	"Dalvik/1.4.0+(Linux;+U;+Android+2.3.5;+UNT+988+Build/MocorDroid2.3.5)",
	"Dalvik/1.4.0+(Linux;+U;+Android+2.3.6;+AP8998+Build/GRK39F)",
	"Dalvik/1.4.0+(Linux;+U;+Android+2.3.6;+G7_SZH+Build/GRK39F)",
	"Dalvik/1.4.0+(Linux;+U;+Android+2.3.6;+IUSAI+i502+Build/GRK39F)",
	"Dalvik/1.4.0+(Linux;+U;+Android+2.3.6;+v60+Build/GRK39F)",
	"Dalvik/1.4.0+(Linux;+U;+Android+4.0.4;+G66+Build/GRK39F)",
	"Dalvik/1.4.0+(Linux;+U;+Android+4.0.4;+M2+Build/MocorDroid2.3.5)",
	"Dalvik/1.4.0+(Linux;+U;+Android+4.0.4;+M3+Build/GRK39F)",
	"Dalvik/1.4.0+(Linux;+U;+Android+4.0.4;+SUNVAN+Build/MocorDroid2.3.5)",
	"Dalvik/1.4.0+(Linux;+U;+Android+4.0.7;+V11+Build/MocorDroid2.3.5)",
	"Dalvik/1.4.0+(Linux;+U;+Android+4.1.1;+MCT999+Build/GRK39F)",
	"Dalvik/1.4.0+(Linux;+U;+Android+5.3.4(Android+2.3.6);+MC919LL+Build/GRK39F)",
	"Dalvik/1.4.0+(Linux;+U;+Android+unknown;+SOP-W168+Build/GRK39F)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.0.3+Build/06.11.001.120626.7728)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.0.4;+BF_A18+Build/IMM76D)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.0.4;+DW+E6+Build/IMM76D)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.0.4;+HYF938A+Build/IMM76D)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.0.4;+T-smart+D68X+Build/alpsD68X)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.0.4;+aigo-P728+Build/IMM76D)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.0.4;+i88+Build/IMM76D)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.1.1;+GT-N7100+Build/IMM76D)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.1.1;+OUKI+A10??+Build/IMM76D)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.1.1;+QX2000+Build/JRO03C)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.1.1;+U8+Build/JRO03C)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.1.1;+newish_L19+Build/JRO03C)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.1.2;+GT-N8000+Build/JZO54K)",
	"Dalvik/1.1.0+(Linux;+U;+Android+2.1-update1;+U20i+Build/2.1.1.A.0.6)",
	"Dalvik/1.4.0+(Linux;+U;+Android+2.3.4;+IMO880+Build/GRJ22)",
	"Dalvik/1.4.0+(Linux;+U;+Android+2.3.5;+ALCATEL+OT+919+Build/GRJ90)",
	"Dalvik/1.4.0+(Linux;+U;+Android+2.3.6;+G3+Build/GRK39F)",
	"Dalvik/1.4.0+(Linux;+U;+Android+2.3.7;+T3696+Build/MocorDroid2.3.7)",
	"Dalvik/1.4.0+(Linux;+U;+Android+4.0.1;+Newish_R610+Build/MocorDroid2.3.5)",
	"Dalvik/1.4.0+(Linux;+U;+Android+4.0.8;+M12++Build/MocorDroid2.3.5)",
	"Dalvik/1.4.0+(Linux;+U;+Android+4.0;+V1268+Build/GRK39F)",
	"Dalvik/1.4.0+(Linux;+U;+Android+4.1.0;+L98+Build/MocorDroid2.3.5)",
	"Dalvik/1.4.0+(Linux;+U;+Android+4.1.9;+GT_I9300+Build/GINGERBREAD.XXKL3)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.0.3;+ETON+T610+Build/IML74K)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.0.4;+AOLE+627+Build/IMM76D)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.0.4;+AOLE+G3+Build/IMM76D)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.0.4;+B063+Build/IMM76D)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.0.4;+ChanghongV9+Build/IMM76D)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.0.4;+I5015C+Build/IMM76D)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.0.4;+KliTON+P188HY+Build/IMM76D)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.0.4;+MT887+Build/6.7.2_GC-125-YTZTD-51)",
	"Dalvik/1.6.0+(Linux;+U;+Android+4.0.4;+Philips+W536+Build/IMM76D)",
	"HuaweiC8650/C92B839+CORE/6.506.4.1+OpenCORE/2.02+(Linux;Android+2.3.3)",
	"HuaweiC8650/C92B879+CORE/6.506.4.1+OpenCORE/2.02+(Linux;Android+2.3.3)",
	"Mozilla/5.0+(compatible;+MSIE+9.0;+Windows+NT+6.1;+WOW64;+Trident/5.0)",
	"SonyEricssonLT26ii+Build/6.1.A.2.45+stagefright/1.2+(Linux;Android+4.0.4)",
	"Mozilla/5.0 (Linux; U; Android 4.4.4; zh-cn; HTC_D820u Build/KTU84P) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
}

func GetRandUserAgent() string {
	return UA[RandRange(0, len(UA))]
}

func init() {
	rand.Seed(time.Now().Unix())
}
func RandRange(min int, max int) int {

	return min + rand.Int()%(max-min+1)
}

func MD5String(raw string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
}
