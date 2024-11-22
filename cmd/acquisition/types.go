package acquisition

const (
	contentType                 = "acquisition"
	supportedAcquisitionVersion = "0.2.0"
)

type DataModel struct {
	AcquisitionVersion string `yaml:"__acquisition_version__"`
	ArchiveUnit        struct {
		Name                 string `yaml:"name"`
		Type                 string `yaml:"type"`
		Creator              string `yaml:"creator"`
		Description          string `yaml:"description"`
		CopyrightClearance   string `yaml:"copyright-clearance"`
		AccessConsiderations string `yaml:"access-considerations"`
		Deposit              struct {
			Depositor          string `yaml:"depositor"`
			Date               string `yaml:"date"`
			AcquisitionPurpose string `yaml:"acquisition-purpose"`
		} `yaml:"deposit"`
		Handling struct {
			Author string `yaml:"author"`
		} `yaml:"handling"`
	} `yaml:"archive-unit"`

	Files []struct {
		Name        string `yaml:"name"`
		Format      string `yaml:"format"`
		Path        string `yaml:"path"`
		Description string `yaml:"description"`
	} `yaml:"files"`
}
