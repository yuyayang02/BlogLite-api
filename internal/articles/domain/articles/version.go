package articles

import (
	"github.com/yuyayang02/BlogLite-api/internal/common/errors"
	"slices"
)

type Version struct {
	hash string
}

func (v Version) Hash() string {
	return v.hash
}

type VersionHistory struct {
	history []Version
}

func (v VersionHistory) RevertTo(versionHash string) (VersionHistory, error) {
	index := slices.IndexFunc(v.history, func(version Version) bool {
		return version.hash == versionHash
	})
	if index < 0 {
		return VersionHistory{}, errors.NewDomainErrorf("该版本不存在：%s", versionHash)
	}

	return VersionHistory{history: v.history[:index]}, nil
}

func (v VersionHistory) Archive(version Version) VersionHistory {
	return VersionHistory{history: append(v.history, version)}
}

func (v VersionHistory) Diff(history VersionHistory) []string {
	var (
		set    = make(map[string]struct{})
		result = make([]string, 0, len(v.history))
	)

	for _, item := range history.history {
		set[item.hash] = struct{}{}
	}

	for _, item := range v.history {
		_, exist := set[item.hash]
		if !exist {
			result = append(result, item.hash)
		}
	}
	return result
}

func newVersionHistory() VersionHistory {
	return VersionHistory{make([]Version, 0)}
}
