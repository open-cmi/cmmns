package prod

func GetProdNav() []Menu {
	return gProdConf.Nav
}

func GetProdInfo() ProdBriefInfo {
	return ProdBriefInfo{
		Name:   gProdConf.Name,
		Footer: gProdConf.Footer,
	}
}
