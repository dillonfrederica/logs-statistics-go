package internal

import (
	"io"
	"net/http"
	"os"
)

func download(url, output string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := os.WriteFile(output, dataBytes, 0644); err != nil {
		return err
	}

	return nil
}

// https://raw.githubusercontent.com/FW27623/qqwry/main/qqwry.dat
func DownloadQqwry(output string) error {
	return download("https://raw.githubusercontent.com/FW27623/qqwry/main/qqwry.dat", output)
}

// https://github.com/P3TERX/GeoLite.mmdb
func DownloadGeoLite(output string) error {
	return download("https://git.io/GeoLite2-City.mmdb", output)
}

// https://github.com/ljxi/GeoCN
func DownloadGeoCN(output string) error {
	return download("https://github.com/ljxi/GeoCN/releases/download/Latest/GeoCN.mmdb", output)
}
