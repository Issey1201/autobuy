package autobuy

type TargetSite interface {
	Run(string) error
	getCheckInfo() map[string]string
}
