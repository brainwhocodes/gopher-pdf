//go:build !cgo || nocgo

package fitz

import (
	"image"
	"io"
)

// Document is a stub when cgo is unavailable.
type Document struct{}

func New(_ string) (*Document, error) {
	return nil, ErrCGODisabled
}

func NewFromMemory(_ []byte) (*Document, error) {
	return nil, ErrCGODisabled
}

func NewFromMemoryWithMagic(_ []byte, _ string) (*Document, error) {
	return nil, ErrCGODisabled
}

func NewFromReader(_ io.Reader) (*Document, error) {
	return nil, ErrCGODisabled
}

func (f *Document) NumPage() int {
	_ = f
	return 0
}

func (f *Document) NeedsPassword() bool {
	_ = f
	return false
}

func (f *Document) AuthenticatePassword(_ string) bool {
	_ = f
	return false
}

func (f *Document) IsReflowable() bool {
	_ = f
	return false
}

func (f *Document) Layout(_, _, _ float64) {
	_ = f
}

func (f *Document) CountChapters() int {
	_ = f
	return 0
}

func (f *Document) CountChapterPages(_ int) int {
	_ = f
	return 0
}

func (f *Document) LastLocation() Location {
	_ = f
	return Location{}
}

func (f *Document) NextLocation(_ Location) Location {
	_ = f
	return Location{}
}

func (f *Document) PreviousLocation(_ Location) Location {
	_ = f
	return Location{}
}

func (f *Document) ClampLocation(_ Location) Location {
	_ = f
	return Location{}
}

func (f *Document) LocationFromPageNumber(_ int) Location {
	_ = f
	return Location{}
}

func (f *Document) PageNumberFromLocation(_ Location) int {
	_ = f
	return 0
}

func (f *Document) MakeBookmark(_ Location) Bookmark {
	_ = f
	return 0
}

func (f *Document) LookupBookmark(_ Bookmark) Location {
	_ = f
	return Location{}
}

func (f *Document) IsDirty() bool {
	_ = f
	return false
}

func (f *Document) CanSaveIncrementally() bool {
	_ = f
	return false
}

func (f *Document) IsFastWebAccess() bool {
	_ = f
	return false
}

func (f *Document) IsRepaired() bool {
	_ = f
	return false
}

func (f *Document) HasPermission(_ int) bool {
	_ = f
	return false
}

func (f *Document) PageLabel(_ int) (string, error) {
	_ = f
	return "", ErrCGODisabled
}

func (f *Document) PageLabelRules() ([]PageLabelRule, error) {
	_ = f
	return nil, ErrCGODisabled
}

func (f *Document) SetPageLabels(_ []PageLabelRule) error {
	_ = f
	return ErrCGODisabled
}

func (f *Document) ResolveLink(_ string) (Location, float64, float64, bool) {
	_ = f
	return Location{}, 0, 0, false
}

func (f *Document) ResolveNames() map[string]NamedDest {
	_ = f
	return map[string]NamedDest{}
}

func (f *Document) XrefLength() int {
	_ = f
	return 0
}

func (f *Document) XrefObject(_ int) (string, error) {
	_ = f
	return "", ErrCGODisabled
}

func (f *Document) XrefGetKey(_ int, _ string) (string, error) {
	_ = f
	return "", ErrCGODisabled
}

func (f *Document) XrefGetKeyTyped(_ int, _ string) (string, string, error) {
	_ = f
	return "", "", ErrCGODisabled
}

func (f *Document) XrefSetKey(_ int, _, _ string) error {
	_ = f
	return ErrCGODisabled
}

func (f *Document) XrefIsFont(_ int) (bool, error) {
	_ = f
	return false, ErrCGODisabled
}

func (f *Document) XrefIsImage(_ int) (bool, error) {
	_ = f
	return false, ErrCGODisabled
}

func (f *Document) XrefIsForm(_ int) (bool, error) {
	_ = f
	return false, ErrCGODisabled
}

func (f *Document) XrefIsStream(_ int) (bool, error) {
	_ = f
	return false, ErrCGODisabled
}

func (f *Document) XrefStream(_ int) ([]byte, error) {
	_ = f
	return nil, ErrCGODisabled
}

func (f *Document) XrefStreamRaw(_ int) ([]byte, error) {
	_ = f
	return nil, ErrCGODisabled
}

func (f *Document) XrefGetKeys(_ int) ([]string, error) {
	_ = f
	return nil, ErrCGODisabled
}

func (f *Document) CatalogXref() int {
	_ = f
	return 0
}

func (f *Document) PageXref(_ int) int {
	_ = f
	return 0
}

func (f *Document) FormFieldCount() int {
	_ = f
	return -1
}

func (f *Document) Image(_ int) (*image.RGBA, error) {
	_ = f
	return nil, ErrCGODisabled
}

func (f *Document) ImageDPI(_ int, _ float64) (*image.RGBA, error) {
	_ = f
	return nil, ErrCGODisabled
}

func (f *Document) ImagePNG(_ int, _ float64) ([]byte, error) {
	_ = f
	return nil, ErrCGODisabled
}

func (f *Document) Links(_ int) ([]Link, error) {
	_ = f
	return nil, ErrCGODisabled
}

func (f *Document) Text(_ int) (string, error) {
	_ = f
	return "", ErrCGODisabled
}

func (f *Document) TextWithFlags(_ int, _ int) (string, error) {
	_ = f
	return "", ErrCGODisabled
}

func (f *Document) HTML(_ int, _ bool) (string, error) {
	_ = f
	return "", ErrCGODisabled
}

func (f *Document) HTMLWithFlags(_ int, _ bool, _ int) (string, error) {
	_ = f
	return "", ErrCGODisabled
}

func (f *Document) XHTML(_ int, _ bool) (string, error) {
	_ = f
	return "", ErrCGODisabled
}

func (f *Document) XHTMLWithFlags(_ int, _ bool, _ int) (string, error) {
	_ = f
	return "", ErrCGODisabled
}

func (f *Document) XML(_ int) (string, error) {
	_ = f
	return "", ErrCGODisabled
}

func (f *Document) XMLWithFlags(_ int, _ int) (string, error) {
	_ = f
	return "", ErrCGODisabled
}

func (f *Document) JSON(_ int, _ float64) (string, error) {
	_ = f
	return "", ErrCGODisabled
}

func (f *Document) JSONWithFlags(_ int, _ int, _ float64) (string, error) {
	_ = f
	return "", ErrCGODisabled
}

func (f *Document) SearchCount(_ int, _ string, _ int, _ int) (int, error) {
	_ = f
	return 0, ErrCGODisabled
}

func (f *Document) Search(_ int, _ string, _ int, _ int, _ *Rect) ([]Rect, error) {
	_ = f
	return nil, ErrCGODisabled
}

func (f *Document) TextBox(_ int, _ Rect, _ int) (string, error) {
	_ = f
	return "", ErrCGODisabled
}

func (f *Document) SVG(_ int) (string, error) {
	_ = f
	return "", ErrCGODisabled
}

func (f *Document) ToC() ([]Outline, error) {
	_ = f
	return nil, ErrCGODisabled
}

func (f *Document) Metadata() map[string]string {
	_ = f
	return map[string]string{}
}

func (f *Document) SetMetadata(_ map[string]string) {
	_ = f
}

func (f *Document) Bound(_ int) (image.Rectangle, error) {
	_ = f
	return image.Rectangle{}, ErrCGODisabled
}

func (f *Document) BoundBox(_ int, _ int) (image.Rectangle, error) {
	_ = f
	return image.Rectangle{}, ErrCGODisabled
}

func (f *Document) PageHasAnnots(_ int) (bool, error) {
	_ = f
	return false, ErrCGODisabled
}

func (f *Document) Close() error {
	_ = f
	return nil
}
