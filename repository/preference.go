package repository

type Unit string
type Axis string

type Preferences struct {
	ID          int                    `storm:"id" json:"-"`
	Dimensions  DimensionPreferences   `storm:"inline" json:"dimensions"`
	Tolerance   float64                `json:"tolerance"`
	UpdateRate  float64                `json:"updateRate"`
	Orientation OrientationPreferences `storm:"inline" json:"orientation"`
}

type DimensionPreferences struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Units  Unit    `json:"units"`
}

type OrientationPreferences struct {
	Length      Axis `json:"length"`
	Width       Axis `json:"width"`
	InvertPitch bool `json:"invertPitch"`
	InvertRoll  bool `json:"invertRoll"`
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
	updated.ID = PreferencesId
	err := r.db.Save(&updated)
	return updated, err
}
