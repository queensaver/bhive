package sound

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"

	"github.com/queensaver/packages/sound"
)

func RecordSound(mac string, duration int, file string) (*sound.Sound, error) {
	r := &sound.Sound{BhiveId: mac, Duration: duration}
	cmd := exec.Command("arecord", "-f", "S16_LE", "-d", fmt.Sprintf("%d", duration), "-r", "44100", "--device=\"hw:1,0\"", file)
	err := cmd.Run()
	if err != nil {
		r.Error = err.Error()
		return r, err
	}
	f, err := os.ReadFile(file)
	if err != nil {
		r.Error = err.Error()
		return r, err
	}
	r.Sound = base64.StdEncoding.EncodeToString(f)
	return r, nil

}
