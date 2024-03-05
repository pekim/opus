package opus

// #include <opus.h>
import "C"

import (
	"fmt"
)

func Test() {
	fmt.Println(C.OPUS_INTERNAL_ERROR)
}
