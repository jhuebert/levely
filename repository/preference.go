package repository

type Unit string
type Axis string

type Preferences struct {
	ID                     int     `storm:"id" json:"id"`
	Version                int     `json:"version"`
	DimensionLength        float64 `json:"dimensionLength"`
	DimensionWidth         float64 `json:"dimensionWidth"`
	DimensionUnits         Unit    `json:"dimensionUnits"`
	OrientationPitch       Axis    `json:"orientationPitch"`
	OrientationRoll        Axis    `json:"orientationRoll"`
	OrientationInvertPitch bool    `json:"orientationInvertPitch"`
	OrientationInvertRoll  bool    `json:"orientationInvertRoll"`
}

const (
	PreferencesId   int  = 1
	UnitInches      Unit = "in"
	UnitCentimeters Unit = "cm"
	AxisX           Axis = "x"
	AxisY           Axis = "y"
	AxisZ           Axis = "z"
)

func (r *Repository) GetPreferences() (Preferences, error) {
	var entity Preferences
	err := r.db.One("ID", PreferencesId, &entity)
	if err != nil {
		return Preferences{}, err
	}
	return entity, nil
}

func (r *Repository) UpdatePreferences(updated Preferences) (Preferences, error) {
	err := r.db.Save(&updated)
	return updated, err
}
