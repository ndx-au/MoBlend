package runner

import (
	"fmt"
	"image/png"
	"os"
)

func validatePNG(path string, resolution [2]int) error {
	st, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !st.Mode().IsRegular() || st.Size() == 0 {
		return fmt.Errorf("PNG must be a nonempty regular file")
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	config, err := png.DecodeConfig(f)
	if err != nil {
		return fmt.Errorf("PNG header: %w", err)
	}
	if config.Width != resolution[0] || config.Height != resolution[1] {
		return fmt.Errorf("PNG dimensions %dx%d, want %dx%d", config.Width, config.Height, resolution[0], resolution[1])
	}
	if _, err = f.Seek(0, 0); err != nil {
		return err
	}
	if _, err = png.Decode(f); err != nil {
		return fmt.Errorf("PNG data: %w", err)
	}
	return nil
}
