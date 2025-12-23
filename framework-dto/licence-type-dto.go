package frameworkdto

type GetLicenceTypeDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MaxSeats    int    `json:"max_seats"`
}

type LicenceTypeCreateRequestDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	MaxSeats    int    `json:"max_seats"`
}

type LicenceTypeUpdateRequestDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MaxSeats    int    `json:"max_seats"`
}
