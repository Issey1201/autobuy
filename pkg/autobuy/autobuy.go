package autobuy

type TargetSite interface {
	Run() error
	getCheckInfo() map[string]string
}
