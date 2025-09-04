package partnerApiClient

import "github.com/Orange-Health/citadel/common/structures"

type GetPartnersApiResponse struct {
	Count int
	Data  []structures.Partner
}
