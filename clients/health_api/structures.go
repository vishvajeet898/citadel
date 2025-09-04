package healthApiClient

import "github.com/Orange-Health/citadel/common/structures"

type DoctorApiResponse struct {
	Result []structures.Doctor `json:"result"`
	Count  int                 `json:"count"`
}
