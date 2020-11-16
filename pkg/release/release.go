package release

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/benjlevesque/ghr/pkg/config"
	"github.com/benjlevesque/ghr/pkg/gh"
	"github.com/benjlevesque/ghr/pkg/util"
	"github.com/google/go-github/v32/github"
	"gopkg.in/cheggaaa/pb.v1"
)

type ReleaseManager struct {
	Owner string
	Repo  string
	Tag   string
}

// Install tries to download, install and save the configuration of the given repo asset.
// If tag is empty, it will default to latest
// If assetName is empty, it will try to match an asset given the current system architecture and OS
func (rm *ReleaseManager) Install(assetName, installPath string) error {
	release, err := gh.GetReleaseByTag(rm.Owner, rm.Repo, rm.Tag)
	if err != nil {
		return err
	}

	asset := getAsset(assetName, release.Assets)

	if asset == nil {
		if assetName != "" {
			return fmt.Errorf("Asset %s not found", assetName)
		}
		return fmt.Errorf("Could not find asset for your system, you can specify the asset to install as a parameter")
	}
	assetName = *asset.Name
	checksum := getChecksum(assetName, release.Assets)
	if checksum == "" {
		fmt.Println("Warning: could not find the checksum")
	}

	executablePath, err := downloadAndInstallAsset(assetName, *asset.BrowserDownloadURL, checksum, installPath)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully installed %s/%s, version %s\n", rm.Owner, rm.Repo, *release.TagName)

	return config.AddOrUpdate(config.ConfigItem{
		Name:       rm.Owner + "/" + rm.Repo,
		Version:    *release.TagName,
		Checksum:   checksum,
		Executable: executablePath,
	})
}

func downloadAndInstallAsset(name, url, checksum, path string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", fmt.Errorf("Cannot get %s: %s", url, err)
	}
	defer resp.Body.Close()
	i, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	sourceSize := int64(i)
	bar := pb.New(int(sourceSize)).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	bar.Start()
	reader := bar.NewProxyReader(resp.Body)
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("Error reading body: %s", err)
	}
	bar.Finish()

	if checksum != "" && checksum != fmt.Sprintf("%x", sha256.Sum256(body)) {
		return "", fmt.Errorf("Checksums don't match")
	}

	installedPath, err := util.ExtractTarGzBinary(bytes.NewReader(body), path)
	if err != nil {
		return "", err
	}
	return installedPath, nil
}

func getArchAliases(arch string) []string {
	switch arch {
	case "386":
		return []string{arch, "i386"}
	case "amd64":
		return []string{arch, "x86_64", "64bit"}
	}
	return []string{arch}
}

func getAsset(assetName string, assets []*github.ReleaseAsset) *github.ReleaseAsset {
	for _, asset := range assets {
		name := strings.ToLower(*asset.Name)
		assetName := strings.ToLower(assetName)
		if assetName != "" && name == assetName {
			return asset
		} else if assetName == "" {
			if !strings.Contains(name, ".tar.gz") {
				continue
			}
			if !strings.Contains(name, runtime.GOOS) {
				continue
			}
			for _, arch := range getArchAliases(runtime.GOARCH) {
				if strings.Contains(name, arch) {
					return asset
				}
			}
		}
	}

	return nil
}

func getChecksum(assetName string, assets []*github.ReleaseAsset) string {
	// look for & download checksums.txt
	for _, asset := range assets {
		if strings.Contains(strings.ToLower(*asset.Name), "checksum") {
			resp, err := http.Get(*asset.BrowserDownloadURL)
			if err == nil {
				body, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					for _, line := range strings.Split(string(body), "\n") {
						// look for the asset name in the file
						// the first line part is the checksum
						parts := strings.Fields(line)
						if len(parts) < 2 {
							continue
						}
						if parts[1] == assetName {
							return parts[0]
						}
					}
				}
			}
		}
	}
	return ""
}
