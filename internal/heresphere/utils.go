package heresphere

import (
	"fmt"

	"github.com/stashapp/stash/internal/manager/config"
)

/*
 * Finds the selected VR Tag string
 */
func getVrTag() (varTag string, err error) {
	// Find setting
	varTag = config.GetInstance().GetUIVRTag()
	if len(varTag) == 0 {
		err = fmt.Errorf("zero length vr tag")
	}
	return
}

/*
 * Finds the selected minimum play percentage value
 */
func getMinPlayPercent() (per int, err error) {
	per = config.GetInstance().GetUIMinPlayPercent()
	if per < 0 {
		err = fmt.Errorf("unset minimum play percent")
	}
	return
}
