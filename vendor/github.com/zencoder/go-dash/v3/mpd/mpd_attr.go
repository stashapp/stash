package mpd

type AttrMPD interface {
	GetStrptr() *string
}

type attrAvailabilityStartTime struct {
	strptr *string
}

func (attr *attrAvailabilityStartTime) GetStrptr() *string {
	return attr.strptr
}

// AttrAvailabilityStartTime returns AttrMPD object for NewMPD
func AttrAvailabilityStartTime(value string) AttrMPD {
	return &attrAvailabilityStartTime{strptr: &value}
}

type attrMinimumUpdatePeriod struct {
	strptr *string
}

func (attr *attrMinimumUpdatePeriod) GetStrptr() *string {
	return attr.strptr
}

// AttrMinimumUpdatePeriod returns AttrMPD object for NewMPD
func AttrMinimumUpdatePeriod(value string) AttrMPD {
	return &attrMinimumUpdatePeriod{strptr: &value}
}

type attrMediaPresentationDuration struct {
	strptr *string
}

func (attr *attrMediaPresentationDuration) GetStrptr() *string {
	return attr.strptr
}

// AttrMediaPresentationDuration returns AttrMPD object for NewMPD
func AttrMediaPresentationDuration(value string) AttrMPD {
	return &attrMediaPresentationDuration{strptr: &value}
}

type attrPublishTime struct {
	strptr *string
}

func (attr *attrPublishTime) GetStrptr() *string {
	return attr.strptr
}

// AttrPublishTime returns AttrMPD object for NewMPD
func AttrPublishTime(value string) AttrMPD {
	return &attrPublishTime{strptr: &value}
}
