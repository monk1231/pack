// +build acceptance

package config

import (
	"strings"

	"github.com/Masterminds/semver"

	"github.com/buildpacks/pack/internal/builder"
)

type LifecycleAsset struct {
	path       string
	descriptor builder.LifecycleDescriptor
}

func (a AssetManager) NewLifecycleAsset(kind ComboValue) LifecycleAsset {
	return LifecycleAsset{
		path:       a.LifecyclePath(kind),
		descriptor: a.LifecycleDescriptor(kind),
	}
}

func (l *LifecycleAsset) Version() string {
	return l.SemVer().String()
}

func (l *LifecycleAsset) SemVer() *builder.Version {
	return l.descriptor.Info.Version
}

func (l *LifecycleAsset) Identifier() string {
	if l.HasLocation() {
		return l.path
	} else {
		return l.Version()
	}
}

func (l *LifecycleAsset) HasLocation() bool {
	return l.path != ""
}

func (l *LifecycleAsset) EscapedPath() string {
	return strings.ReplaceAll(l.path, `\`, `\\`)
}

func (l *LifecycleAsset) BuildpackAPIVersion() string {
	return l.descriptor.API.BuildpackVersion.String()
}

func (l *LifecycleAsset) PlatformAPIVersion() string {
	return l.descriptor.API.PlatformVersion.String()
}

func (l *LifecycleAsset) ShouldShowReference() bool {
	return !l.SemVer().LessThan(semver.MustParse("0.5.0"))
}

func (l *LifecycleAsset) ShouldShowProcesses() bool {
	return !l.pre060()
}

type LifecycleFeature int

const (
	DefaultProcess LifecycleFeature = iota
	CreatorInLifecycle
	DetailedCacheLogging
)

var lifecycleFeatureTests = map[LifecycleFeature]func(l *LifecycleAsset) bool{
	DefaultProcess: func(l *LifecycleAsset) bool {
		return l.atLeast070()
	},
	CreatorInLifecycle: func(l *LifecycleAsset) bool {
		return l.atLeast074()
	},
	DetailedCacheLogging: func(l *LifecycleAsset) bool {
		return !l.pre060()
	},
}

func (l *LifecycleAsset) SupportsFeature(f LifecycleFeature) bool {
	return lifecycleFeatureTests[f](l)
}

func (l *LifecycleAsset) pre060() bool {
	return l.SemVer().LessThan(semver.MustParse("0.6.0"))
}

func (l *LifecycleAsset) atLeast070() bool {
	return !l.SemVer().LessThan(semver.MustParse("0.7.0"))
}

func (l *LifecycleAsset) atLeast074() bool {
	return !l.SemVer().LessThan(semver.MustParse("0.7.4"))
}
