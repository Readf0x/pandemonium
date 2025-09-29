package systems

import (
	"fmt"
	"math"
	"time"
)

func TimeAgo(t time.Time) string {
	d := time.Since(t)

	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	case d < 7*24*time.Hour:
		return fmt.Sprintf("%dd", int(d.Hours())/24)
	default:
		return t.Format("2006-01-02")
	}
}

func Trunc32(n int32) string {
	numstr := fmt.Sprintf("%d", n)
	reps := []string{"K", "M", "B"}
	for i := 3; i < 3*len(reps); i += 3 {
		if len(numstr) > i {
			return fmt.Sprintf("%.1f%s", float64(n)/math.Pow(1000, float64(i/3)), reps[i/3-1])
		}
	}
	return numstr
}
func Trunc64(n int64) string {
	numstr := fmt.Sprintf("%d", n)
	reps := []string{"K", "M", "B"}
	for i := 3; i < 3*len(reps); i += 3 {
		if len(numstr) > i {
			return fmt.Sprintf("%.1f%s", float64(n)/math.Pow(1000, float64(i/3)), reps[i/3-1])
		}
	}
	return numstr
}

func Encode(source string) (fin string) {
	for _, r := range []rune(source) {
		fmt.Printf("%d, %c\n",r,r)
		fin += fmt.Sprintf("\\u%04X", r)
	}
	return
}
